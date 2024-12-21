package pilot

import (
	"context"
	baseconfig "github.com/skiff-sh/config"
	"github.com/skiff-sh/config/ptr"
	"github.com/skiff-sh/pilot/server/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var (
	DefaultName              = "pilot"
	DefaultNamespace         = "pilot"
	DefaultContainerPort     = int32(8080)
	DefaultContainerPortName = "main"
	grpcProbe                = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			GRPC: &corev1.GRPCAction{
				Port: DefaultContainerPort,
			},
		},
	}
)

func Deploy(ctx context.Context, cl kubernetes.Interface, image string, conf *config.Config) error {
	evs := baseconfig.ToEnvVars("pilot", conf)
	ing := newIngress(DefaultName, DefaultNamespace)
	svc := newService(DefaultName, DefaultNamespace)
	dep := newDeployment(DefaultName, DefaultNamespace, image, evs)
	co := metav1.CreateOptions{}
	var err error
	dep, err = cl.AppsV1().Deployments(DefaultNamespace).Create(ctx, dep, co)
	if err != nil {
		return err
	}

	svc, err = cl.CoreV1().Services(DefaultNamespace).Create(ctx, svc, co)
	if err != nil {
		return err
	}

	ing, err = cl.NetworkingV1().Ingresses(DefaultNamespace).Create(ctx, ing, co)
	if err != nil {
		return err
	}

	return nil
}

func newIngress(name, namespace string) *networkingv1.Ingress {
	out := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.Ptr(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: name,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
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
					TargetPort: intstr.FromString(DefaultContainerPortName),
				},
			},
			Selector: map[string]string{
				"app": name,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	return out
}

func newDeployment(name, namespace, image string, envVars map[string]string) *appsv1.Deployment {
	evs := make([]corev1.EnvVar, 0, len(envVars))
	for k, v := range envVars {
		evs = append(evs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

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
							Name:  "main",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          DefaultContainerPortName,
									ContainerPort: DefaultContainerPort,
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
