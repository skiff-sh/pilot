package pilot

import (
	"context"
	"github.com/skiff-sh/ksuite"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestPilotSuite struct {
	ksuite.KubeSuite
}

func (t *TestPilotSuite) TestDeploy() {
	// -- Given
	//
	ctx := context.TODO()
	Deploy(ctx, t.Kube)

	// -- When
	//

	// -- Then
	//
}

func TestTestPilotSuite(t *testing.T) {
	suite.Run(t, new(TestPilotSuite))
}
