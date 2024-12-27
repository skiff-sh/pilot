package pilot

import (
	"context"
	"fmt"
	"github.com/skiff-sh/pilot/server/pkg/pilot"
	"net/http"
	"testing"
	"time"

	"github.com/skiff-sh/ksuite"
	"github.com/skiff-sh/kube"
	"github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/labels"
)

type TestPilotSuite struct {
	ksuite.KubeSuite
}

func (t *TestPilotSuite) TestE2E() {
	type deps struct {
		Cl  pilot.Client
		Ctx context.Context
	}

	type test struct {
		ProvokeName  string
		Constructor  func(d *deps) test
		Conduct      func(ctx context.Context, cl pilot.Client) error
		ExpectedFunc func()
		ExpectedErr  string
	}

	hit := make(chan string, 1)
	go func() {
		http.HandleFunc("/derp", func(writer http.ResponseWriter, request *http.Request) {
			hit <- request.Header.Get("hi")
		})

		server := &http.Server{
			Addr:              ":8085",
			ReadHeaderTimeout: 5 * time.Second,
		}
		_ = server.ListenAndServe()
	}()

	tests := map[string]test{
		"http req": {
			ProvokeName: "request",
			ExpectedFunc: func() {
				t.Equal("there", ExpectWithin(&t.Suite, hit, 5*time.Second))
			},
			Conduct: func(ctx context.Context, cl pilot.Client) error {
				_, err := cl.NewBehavior().Name("request").
					Tendency().
					Action().
					HTTPRequest(fmt.Sprintf("http://%s:8085/derp", t.Cluster.HostIP), pilot.WithHTTPHeader("hi", "there")).
					Add().
					Send(ctx)
				return err
			},
		},
	}

	conf, err := config.New()
	if !t.NoError(err) {
		return
	}
	podCl := t.Kube.CoreV1().Pods(conf.Test.Namespace)
	selector := labels.SelectorFromValidatedSet(labels.Set{"app": conf.Test.DeployName})
	for desc, v := range tests {
		t.Run(desc, func() {
			ctx := t.Cluster.Context

			d := &deps{}
			if v.Constructor != nil {
				v = v.Constructor(d)
			}

			err = pilot.DeployK8s(ctx, t.Kube, conf)
			if !t.NoError(err) {
				return
			}

			_, err = kube.WaitPodReady(ctx, podCl, selector)
			if !t.NoError(err) {
				return
			}

			cl, err := pilot.Connect(ctx, fmt.Sprintf("localhost:%d", t.Cluster.ExposedNodePort))
			if !t.NoError(err) {
				return
			}

			err = v.Conduct(ctx, cl)
			if !t.NoError(err) {
				return
			}

			_, err = cl.Provoke(ctx, v.ProvokeName)
			if v.ExpectedErr != "" || !t.NoError(err) {
				t.EqualError(err, v.ExpectedErr)
				return
			}

			if v.ExpectedFunc != nil {
				v.ExpectedFunc()
			}
		})
	}
}

func ExpectWithin[T any](t *suite.Suite, c chan T, to time.Duration) (out T) {
	timer := time.NewTimer(to)
	select {
	case <-timer.C:
		t.Fail("took too long to receive request")
		return
	case hit := <-c:
		return hit
	}
}

func TestTestPilotSuite(t *testing.T) {
	suite.Run(t, new(TestPilotSuite))
}
