// +build integration

package sdk_test

import (
	"context"
	"testing"

	"github.com/src-d/lookout/util/cmdtest"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	ctx  context.Context
	stop func()
}

func (suite *IntegrationSuite) SetupSuite() {
	suite.ctx, suite.stop = cmdtest.StoppableCtx()
	cmdtest.StartDummy(suite.ctx, suite.Require())
}

func (suite *IntegrationSuite) TearDownSuite() {
	suite.stop()
}

func (suite *IntegrationSuite) TestReview() {
	r := cmdtest.RunCli(suite.ctx, suite.Require(), "review", "ipv4://localhost:10302")
	cmdtest.GrepTrue(r, "posting analysis")
}

func (suite *IntegrationSuite) TestPush() {
	r := cmdtest.RunCli(suite.ctx, suite.Require(), "push", "ipv4://localhost:10302")
	cmdtest.GrepTrue(r, "dummy comment for push event")
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
