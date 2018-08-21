// +build integration

package server_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/src-d/lookout/util/cmdtest"

	"github.com/stretchr/testify/suite"
)

type MultiDummyIntegrationSuite struct {
	suite.Suite
	ctx  context.Context
	stop func()
	r    io.Reader
	w    io.WriteCloser
}

func (suite *MultiDummyIntegrationSuite) SetupTest() {
	cmdtest.ResetDB()

	suite.ctx, suite.stop = cmdtest.StoppableCtx()
	cmdtest.StartDummy(suite.ctx)
	cmdtest.StartDummy(suite.ctx, "--analyzer", "ipv4://localhost:10303")
	suite.r, suite.w = cmdtest.StartServe(suite.ctx, "--provider", "json", "-c",
		"../../fixtures/double_dummy_config.yml", "dummy-repo-url")

	// make sure server started correctly
	cmdtest.GrepTrue(suite.r, "Starting watcher")
}

func (suite *MultiDummyIntegrationSuite) TearDownTest() {
	suite.stop()
}

func (suite *MultiDummyIntegrationSuite) sendEvent(json string) {
	_, err := fmt.Fprintln(suite.w, json)
	suite.Require().NoError(err)
}

func (suite *MultiDummyIntegrationSuite) TestSuccessReview() {
	suite.sendEvent(successJSON)
	cmdtest.GrepTrue(suite.r, "processing pull request")
	cmdtest.GrepTrue(suite.r, "posting analysis")
	found, buf := cmdtest.Grep(suite.r, `status=success`)

	suite.Require().Truef(found, "'%s' not found in:\n%s", `status=success`, buf.String())

	// TODO (carlosms): bug, buf.String() returns the unread portion of the buffer
	// only, test may fail if results come in different order
	suite.Require().Contains(
		buf.String(),
		`{"analyzer-name":"Dummy1","file":"provider/common.go","text":"The file has increased in 5 lines."}`,
		"no comments from the first analyzer")

	suite.Require().Contains(
		buf.String(),
		`{"analyzer-name":"Dummy2","file":"provider/common.go","text":"The file has increased in 5 lines."}`,
		"no comments from the second analyzer")
}

func TestMultiDummyIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MultiDummyIntegrationSuite))
}
