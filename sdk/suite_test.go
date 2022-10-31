package sdk

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite

	client *Client
}

func (ts *testSuite) SetupSuite() {
	client, err := NewClient(nil)
	if err != nil {
		ts.T().Fatalf("new client: %v", err)
	}
	ts.client = client
}

func (ts *testSuite) TearDownSuite() {
	if ts.client != nil {
		ts.client.Close()
	}
}
