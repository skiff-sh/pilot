package pilot

import (
	"context"
	"fmt"
	"github.com/skiff-sh/ksuite"
	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"testing"
)

type TestPilotSuite struct {
	ksuite.KubeSuite
}

func (t *TestPilotSuite) TestDeploy() {
	// -- Given
	//
	ctx := context.TODO()
	err := Deploy(ctx, t.Kube, &config.Config{})
	if !t.NoError(err) {
		return
	}
	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", t.Cluster.IngressPort))
	if !t.NoError(err) {
		return
	}

	cl := New(pilot.NewPilotServiceClient(cc))

	// -- When
	//
	_, err = cl.NewBehavior().Name("derp").Tendency().Action().Exec("echo", WithExecArgs("derp")).Add().Send(ctx)

	// -- Then
	//
	t.NoError(err)
}

func TestTestPilotSuite(t *testing.T) {
	suite.Run(t, new(TestPilotSuite))
}
