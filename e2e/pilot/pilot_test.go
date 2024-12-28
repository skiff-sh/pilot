package pilot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/skiff-sh/pilot/pkg/testutil"
	"github.com/skiff-sh/pilot/server/pkg/pilot"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
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

func (t *TestPilotSuite) SetupSuite() {
	t.SkipCleanNamespaces = []string{"pilot"}
	t.KubeSuite.SetupSuite()
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
		ExpectedFunc func(out *structpb.Struct)
		ExpectedErr  string
	}

	httpHit := make(chan *http.Request, 1)
	go func() {
		http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
			cl := request.Clone(request.Context())
			body, _ := io.ReadAll(request.Body)
			cl.Body = io.NopCloser(bytes.NewReader(body))
			httpHit <- cl
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
			ExpectedFunc: func(_ *structpb.Struct) {
				req := testutil.ExpectWithin(&t.Suite, httpHit, 5*time.Second)
				if !t.NotNil(req) {
					return
				}

				body, _ := io.ReadAll(req.Body)
				t.Equal("there", req.Header.Get("hi"))
				t.Equal("thing", req.Header.Get("some"))
				t.Equal(`{"hello":"hi"}`, string(body))
				t.Equal(http.MethodPost, req.Method)
			},
			Conduct: func(ctx context.Context, cl pilot.Client) error {
				_, err := cl.NewBehavior().Name("request").
					Tendency().
					Action().
					HTTPRequest(fmt.Sprintf("http://%s:8085/test", t.Cluster.HostIP),
						pilot.WithHTTPHeaders(map[string]string{"some": "thing"}),
						pilot.WithHTTPHeader("hi", "there"),
						pilot.WithHTTPJSONBody(map[string]string{"hello": "hi"}),
						pilot.WithHTTPMethod(http.MethodPost),
					).
					Send(ctx)
				return err
			},
		},
		"exec req": {
			ProvokeName: "exec",
			ExpectedFunc: func(out *structpb.Struct) {
				t.Equal("derp\n/etc\n", out.Fields["field"].GetStringValue())
			},
			Conduct: func(ctx context.Context, cl pilot.Client) error {
				_, err := cl.NewBehavior().Name("exec").
					Tendency().
					ID("exec").
					Action().
					Exec("/bin/sh", pilot.WithExecArgs("-c", "echo derp && pwd"), pilot.WithExecDir("/etc")).
					Tendency().
					Action().
					SetResponseField("exec.stdout", "field").
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
	ctx := t.Cluster.Context
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
	for desc, v := range tests {
		t.Run(desc, func() {
			d := &deps{}
			if v.Constructor != nil {
				v = v.Constructor(d)
			}

			err = v.Conduct(ctx, cl)
			if !t.NoError(err) {
				return
			}

			out, err := cl.Provoke(ctx, v.ProvokeName)
			if v.ExpectedErr != "" || !t.NoError(err) {
				t.EqualError(err, v.ExpectedErr)
				return
			}

			if v.ExpectedFunc != nil {
				v.ExpectedFunc(out)
			}
		})
	}
}

func TestTestPilotSuite(t *testing.T) {
	suite.Run(t, new(TestPilotSuite))
}
