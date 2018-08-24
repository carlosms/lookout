// +build integration

package server_test

import (
	"context"
	"fmt"
	"io"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	ctx  context.Context
	stop func()
	r    io.Reader
	w    io.WriteCloser
}

func (suite *IntegrationSuite) sendEvent(json string) {
	_, err := fmt.Fprintln(suite.w, json)
	suite.Require().NoError(err)
}
