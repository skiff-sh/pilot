package pilot

import (
	"context"
	"github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

type DeployK8STestSuite struct {
	suite.Suite
}

func (d *DeployK8STestSuite) TestConstructors() {
	// -- Given
	//
	cl := fake.NewClientset()
	ctx := context.TODO()
	conf, _ := config.New()
	getOpts := metav1.GetOptions{}

	// -- When
	//
	err := DeployK8s(ctx, cl, conf, WithExposeHTTPNodePort(1), WithExposeGRPCNodePort(1))

	// -- Then
	//
	if d.NoError(err) {
		_, err = cl.AppsV1().Deployments(conf.Test.Namespace).Get(ctx, conf.Test.DeployName, getOpts)
		d.NoError(err)
		_, err = cl.CoreV1().Services(conf.Test.Namespace).Get(ctx, conf.Test.DeployName, getOpts)
		d.NoError(err)
		_, err = cl.CoreV1().Namespaces().Get(ctx, conf.Test.DeployName, getOpts)
		d.NoError(err)
	}
}

func TestDeployK8STestSuite(t *testing.T) {
	suite.Run(t, new(DeployK8STestSuite))
}
