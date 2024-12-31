package pilot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skiff-sh/pilot/api/go/pilot"

	baseconfig "github.com/skiff-sh/config"
	"github.com/skiff-sh/config/ptr"
	"github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/skiff-sh/serverapp"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var (
	DefaultGRPCPortName  = "grpc"
	DefaultHTTPPortName  = "http"
	DefaultContainerName = "main"
	DefaultGRPCPort      = int32(81)
	DefaultHTTPPort      = int32(80)
)

func Connect(ctx context.Context, grpcAddr, httpAddr string) (Client, error) {
	cc, err := grpc.NewClient(grpcAddr, serverapp.DefaultDialOpts()...)
	if err != nil {
		return nil, err
	}

	err = serverapp.WaitUntilReady(ctx, grpcAddr, 30*time.Second)
	if err != nil {
		return nil, err
	}

	httpCl := *http.DefaultClient
	httpCl.Timeout = 5 * time.Second

	return New(pilot.NewPilotServiceClient(cc), NewHTTP(httpAddr, &httpCl)), nil
}

type DeployOpts struct {
	GRPCNodePort int32
	HTTPNodePort int32
}

type DeployOpt func(o *DeployOpts)

func WithExposeGRPCNodePort(port uint16) DeployOpt {
	return func(o *DeployOpts) {
		o.GRPCNodePort = int32(port)
	}
}

func WithExposeHTTPNodePort(port uint16) DeployOpt {
	return func(o *DeployOpts) {
		o.HTTPNodePort = int32(port)
	}
}

func DeployK8s(ctx context.Context, cl kubernetes.Interface, conf *config.Config, o ...DeployOpt) error {
	op := &DeployOpts{}
	for _, v := range o {
		v(op)
	}
	evs := baseconfig.ToEnvVars("pilot", conf)
	svc := newService(conf.Test.DeployName, conf.Test.Namespace, op.GRPCNodePort, op.HTTPNodePort)
	dep := newDeployment(conf.Test.DeployName, conf.Test.Namespace, conf.Test.Image, conf.GRPC.Addr.Port(), conf.HTTP.Addr.Port(), evs)
	ns := newNamespace(conf.Test.DeployName)
	co := metav1.CreateOptions{}
	var err error
	_, err = cl.CoreV1().Namespaces().Create(ctx, ns, co)
	if err != nil {
		return err
	}

	_, err = cl.AppsV1().Deployments(conf.Test.Namespace).Create(ctx, dep, co)
	if err != nil {
		return err
	}

	_, err = cl.CoreV1().Services(conf.Test.Namespace).Create(ctx, svc, co)
	if err != nil {
		return err
	}

	return nil
}

func newNamespace(name string) *corev1.Namespace {
	out := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return out
}

func newService(name, namespace string, grpcNodePort, httpNodePort int32) *corev1.Service {
	out := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       DefaultGRPCPortName,
					Port:       DefaultGRPCPort,
					NodePort:   grpcNodePort,
					TargetPort: intstr.FromString(DefaultGRPCPortName),
				},
				{
					Name:       DefaultHTTPPortName,
					Port:       DefaultHTTPPort,
					NodePort:   httpNodePort,
					TargetPort: intstr.FromString(DefaultHTTPPortName),
				},
			},
			Selector: map[string]string{
				"app": name,
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	return out
}

func newDeployment(name, namespace, image string, grpcPort, httpPort uint16, envVars map[string]string) *appsv1.Deployment {
	evs := make([]corev1.EnvVar, 0, len(envVars))
	for k, v := range envVars {
		evs = append(evs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	grpcProbe := newProbe(grpcPort)

	out := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Replicas: ptr.Ptr[int32](1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  DefaultContainerName,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          DefaultGRPCPortName,
									ContainerPort: int32(grpcPort),
								},
								{
									Name:          DefaultHTTPPortName,
									ContainerPort: int32(httpPort),
								},
							},
							Env:            evs,
							LivenessProbe:  grpcProbe,
							ReadinessProbe: grpcProbe,
						},
					},
				},
			},
		},
	}
	return out
}

func newProbe(port uint16) *corev1.Probe {
	return &corev1.Probe{
		PeriodSeconds: int32(3),
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/grpc_health_probe",
					fmt.Sprintf("-addr=localhost:%d", port),
				},
			},
		},
	}
}
