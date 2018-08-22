package cmdtest

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// CmdTimeout defines timeout for a command
var CmdTimeout = time.Minute

// default path to binaries
var dummyBin = "../../build/bin/dummy"
var lookoutBin = "../../build/bin/lookout"

// function to stop running commands
// redefined in StoppableCtx
var stop func()

func init() {
	if os.Getenv("DUMMY_BIN") != "" {
		dummyBin = os.Getenv("DUMMY_BIN")
	}
	if os.Getenv("LOOKOUT_BIN") != "" {
		lookoutBin = os.Getenv("LOOKOUT_BIN")
	}
}

// StoppableCtx return ctx and stop function
func StoppableCtx() (context.Context, func()) {
	ctx, timeoutCancel := context.WithTimeout(context.Background(), CmdTimeout)

	ctx, cancel := context.WithCancel(ctx)
	stop = func() {
		timeoutCancel()
		cancel()
		fmt.Println("stopping services")
		time.Sleep(time.Second) // go needs a bit of time to kill process
	}

	return ctx, stop
}

// StartDummy starts dummy analyzer with context and optional arguments
func StartDummy(ctx context.Context, require *require.Assertions, args ...string) io.Reader {
	r, outputWriter := io.Pipe()
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)

	args = append([]string{"serve"}, args...)

	cmd := exec.CommandContext(ctx, dummyBin, args...)
	cmd.Stdout = outputWriter
	cmd.Stderr = outputWriter
	err := cmd.Start()
	require.NoError(err, "can't start analyzer")

	go func() {
		if err := cmd.Wait(); err != nil {
			// don't print error if analyzer was killed by cancel
			if ctx.Err() != context.Canceled {
				fmt.Println("analyzer exited with error:", err)
				fmt.Printf("output:\n%s", buf.String())
				// T.Fail cannot be called from a goroutine
				failExit()
			}
		}
	}()

	return tee
}

// StartServe starts lookout server with context and optional arguments
func StartServe(ctx context.Context, require *require.Assertions, args ...string) (io.Reader, io.WriteCloser) {
	r, outputWriter := io.Pipe()
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)

	args = append([]string{"serve"}, args...)

	cmd := exec.CommandContext(ctx, lookoutBin, args...)
	cmd.Stdout = outputWriter
	cmd.Stderr = outputWriter

	w, err := cmd.StdinPipe()
	require.NoError(err, "can't start server")

	err = cmd.Start()
	require.NoError(err, "can't start server")

	go func() {
		if err := cmd.Wait(); err != nil {
			// don't print error if analyzer was killed by cancel
			if ctx.Err() != context.Canceled {
				fmt.Println("server exited with error:", err)
				fmt.Printf("output:\n%s", buf.String())
				// T.Fail cannot be called from a goroutine
				failExit()
			}
		}
	}()

	return tee, w
}

// RunCli runs lookout subcommand (not a server)
func RunCli(ctx context.Context, require *require.Assertions, cmd string, args ...string) io.Reader {
	args = append([]string{cmd}, args...)

	var out bytes.Buffer
	cliCmd := exec.CommandContext(ctx, lookoutBin, args...)
	cliCmd.Stdout = &out
	cliCmd.Stderr = &out

	err := cliCmd.Run()
	require.NoErrorf(err,
		"'lookout %s' command returned error. output:\n%s",
		strings.Join(args, " "), out.String())

	return &out
}

// ResetDB recreates database for the test
func ResetDB(require *require.Assertions) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/lookout?sslmode=disable")
	require.NoError(err, "can't connect to DB")

	_, err = db.Exec("DROP SCHEMA public CASCADE;")
	require.NoError(err, "can't execute query")
	_, err = db.Exec("CREATE SCHEMA public;")
	require.NoError(err, "can't execute query")
	_, err = db.Exec("GRANT ALL ON SCHEMA public TO postgres;")
	require.NoError(err, "can't execute query")
	_, err = db.Exec("GRANT ALL ON SCHEMA public TO public;")
	require.NoError(err, "can't execute query")

	err = exec.Command(lookoutBin, "migrate").Run()
	require.NoError(err, "can't migrate DB")
}

func failExit() {
	stop()

	os.Exit(1)
}
