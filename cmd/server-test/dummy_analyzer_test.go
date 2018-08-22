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

type DummyIntegrationSuite struct {
	suite.Suite
	ctx  context.Context
	stop func()
	r    io.Reader
	w    io.WriteCloser
}

func (suite *DummyIntegrationSuite) SetupTest() {
	cmdtest.ResetDB(suite.Require())

	suite.ctx, suite.stop = cmdtest.StoppableCtx()
	cmdtest.StartDummy(suite.ctx, suite.Require())
	suite.r, suite.w = cmdtest.StartServe(suite.ctx, suite.Require(),
		"--provider", "json", "-c", "../../fixtures/dummy_config.yml", "dummy-repo-url")

	// make sure server started correctly
	cmdtest.GrepTrue(suite.r, "Starting watcher")
}

func (suite *DummyIntegrationSuite) TearDownTest() {
	suite.stop()
}

func (suite *DummyIntegrationSuite) sendEvent(json string) {
	_, err := fmt.Fprintln(suite.w, json)
	suite.Require().NoError(err)
}

const successJSON = `{"event":"review", "internal_id": "1", "number": 1, "commit_revision":{"base":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"4eebef102d7979570aadf69ff54ae1ffcca7ce00"},"head":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"d304499cb2a9cad3ea260f06ad59c1658db4763d"}}}`

func (suite *DummyIntegrationSuite) TestSuccessReview() {
	suite.sendEvent(successJSON)
	cmdtest.GrepTrue(suite.r, "processing pull request")
	cmdtest.GrepTrue(suite.r, `{"analyzer-name":"Dummy","file":"provider/common.go","text":"The file has increased in 5 lines."}`)
	cmdtest.GrepTrue(suite.r, `status=success`)
}

func (suite *DummyIntegrationSuite) TestSkipReview() {
	suite.sendEvent(successJSON)
	cmdtest.GrepTrue(suite.r, `status=success`)

	suite.sendEvent(successJSON)
	cmdtest.GrepTrue(suite.r, `event successfully processed, skipping...`)
}

func (suite *DummyIntegrationSuite) TestReviewDontPost() {
	suite.sendEvent(successJSON)
	cmdtest.GrepTrue(suite.r, `status=success`)

	json := `{"event":"review", "internal_id": "2", "number": 1, "commit_revision":{"base":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"4eebef102d7979570aadf69ff54ae1ffcca7ce00"},"head":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"d304499cb2a9cad3ea260f06ad59c1658db4763d"}}}`
	suite.sendEvent(json)
	cmdtest.GrepTrue(suite.r, "processing pull request")
	cmdtest.GrepAndNot(suite.r, `status=success`, `posting analysis`)
}

func (suite *DummyIntegrationSuite) TestWrongRevision() {
	json := `{"event":"review", "internal_id": "3", "number": 3, "commit_revision": {"base":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"0000000000000000000000000000000000000000"},"head":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"0000000000000000000000000000000000000000"}}}`
	suite.sendEvent(json)
	cmdtest.GrepTrue(suite.r, `event processing failed`)
}

func (suite *DummyIntegrationSuite) TestSuccessPush() {
	successPushJSON := `{"event":"push", "internal_id": "1", "commit_revision":{"base":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"4eebef102d7979570aadf69ff54ae1ffcca7ce00"},"head":{"internal_repository_url":"https://github.com/src-d/lookout.git","reference_name":"refs/heads/master","hash":"d304499cb2a9cad3ea260f06ad59c1658db4763d"}}}`
	suite.sendEvent(successPushJSON)
	cmdtest.GrepTrue(suite.r, "processing push")
	cmdtest.GrepTrue(suite.r, "comments can belong only to review event but 1 is given")
	cmdtest.GrepTrue(suite.r, `status=success`)
}

func TestDummyIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DummyIntegrationSuite))
}
