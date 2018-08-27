// +build integration

package sdk_test

import (
	"testing"

	"github.com/src-d/lookout/util/cmdtest"

	"github.com/stretchr/testify/suite"
)

type SDKDummyTestSuite struct {
	cmdtest.IntegrationSuite
}

func (suite *SDKDummyTestSuite) SetupSuite() {
	suite.StoppableCtx()
	suite.StartDummy()
}

func (suite *SDKDummyTestSuite) TearDownSuite() {
	suite.Stop()
}

func (suite *SDKDummyTestSuite) TestReview() {
	r := suite.RunCli("review", "ipv4://localhost:10302")
	cmdtest.GrepTrue(r, "posting analysis")
}

func (suite *SDKDummyTestSuite) TestPush() {
	r := suite.RunCli("push", "ipv4://localhost:10302")
	cmdtest.GrepTrue(r, "dummy comment for push event")
}

func SDKTestDummyTestSuite(t *testing.T) {
	suite.Run(t, new(SDKDummyTestSuite))
}
