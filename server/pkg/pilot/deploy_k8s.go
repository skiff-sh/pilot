package pilot

import (
	"context"
	"fmt"
	"time"

	baseconfig "github.com/skiff-sh/config"
	"github.com/skiff-sh/config/ptr"
	"github.com/skiff-sh/ksuite"
	pilot "github.com/skiff-sh/pilot/api/go"
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
	DefaultContainerPortName = "main"
	DefaultContainerName     = "main"
)

func Connect(ctx context.Context, addr string) (Client, error) {
	cc, err := grpc.NewClient(addr, serverapp.DefaultDialOpts()...)
	if err != nil {
		return nil, err
	}

	err = serverapp.WaitUntilReady(ctx, addr, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return New(pilot.NewPilotServiceClient(cc)), nil
}

func DeployK8s(ctx context.Context, cl kubernetes.Interface, conf *config.Config) error {
	evs := baseconfig.ToEnvVars("pilot", conf)
	svc := newService(conf.Test.DeployName, conf.Test.Namespace)
	dep := newDeployment(conf.Test.DeployName, conf.Test.Namespace, conf.Test.Image, conf.Server.Addr.Port(), evs)
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

func newService(name, namespace string) *corev1.Service {
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
					Name:       DefaultContainerPortName,
					Port:       80,
					NodePort:   int32(ksuite.InternalNodePort),
					TargetPort: intstr.FromString(DefaultContainerPortName),
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

func newDeployment(name, namespace, image string, containerPort uint16, envVars map[string]string) *appsv1.Deployment {
	evs := make([]corev1.EnvVar, 0, len(envVars))
	for k, v := range envVars {
		evs = append(evs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	probe := newProbe(containerPort)

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
									Name:          DefaultContainerPortName,
									ContainerPort: int32(containerPort),
								},
							},
							Env:            evs,
							LivenessProbe:  probe,
							ReadinessProbe: probe,
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
