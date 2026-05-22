package handlers

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var k8sClient *kubernetes.Clientset

func InitK8s() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	k8sClient, err = kubernetes.NewForConfig(config)
	return err
}

func getNamespace() string {
	return getEnvK8s("POD_NAMESPACE", "vm-manager")
}

func getEnvK8s(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func createK8sSecret(name, password string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: getNamespace(),
			Labels: map[string]string{
				"app": "vm-manager-hypervisor",
			},
		},
		StringData: map[string]string{
			"password": password,
		},
	}
	_, err := k8sClient.CoreV1().Secrets(getNamespace()).Create(
		context.Background(), secret, metav1.CreateOptions{},
	)
	return err
}

func deleteK8sSecret(name string) error {
	return k8sClient.CoreV1().Secrets(getNamespace()).Delete(
		context.Background(), name, metav1.DeleteOptions{},
	)
}
