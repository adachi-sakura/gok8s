package appclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8SClient kubernetes.Clientset

func NewK8SClient(config *rest.Config) *K8SClient {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return (*K8SClient)(clientSet)
}

func (cli *K8SClient) K8SClient() *kubernetes.Clientset {
	return (*kubernetes.Clientset)(cli)
}

func (cli *K8SClient) ListPods(namespace string) (*corev1.PodList, error) {
	return cli.K8SClient().CoreV1().Pods(namespace).List(metav1.ListOptions{})
}

func (cli *K8SClient) ListNodes() (*corev1.NodeList, error) {
	return cli.K8SClient().CoreV1().Nodes().List(metav1.ListOptions{})
}

func (cli *K8SClient) GetNode(name string) (*corev1.Node, error) {
	return cli.K8SClient().CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (cli *K8SClient) UpdateNode(node *corev1.Node) (*corev1.Node, error) {
	return cli.K8SClient().CoreV1().Nodes().Update(node)
}

func (cli *K8SClient) CreateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	namespace := deployment.Namespace
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	return cli.K8SClient().AppsV1().Deployments(namespace).Create(deployment)
}

func (cli *K8SClient) CreateDeployments(deployments []*appsv1.Deployment) ([]*appsv1.Deployment, error) {
	createdDeployments := []*appsv1.Deployment{}
	for _, deployment := range deployments {
		createdDeployment, err := cli.CreateDeployment(deployment)
		if err != nil {
			return createdDeployments, err
		}
		createdDeployments = append(createdDeployments, createdDeployment)
	}
	return createdDeployments, nil
}

func (cli *K8SClient) ListLimitRange(namespace string) (*corev1.LimitRangeList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceDefault
	}
	return cli.K8SClient().CoreV1().LimitRanges(namespace).List(metav1.ListOptions{})
}

func (cli *K8SClient) ListResourceQuota(namespace string) (*corev1.ResourceQuotaList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceDefault
	}
	return cli.K8SClient().CoreV1().ResourceQuotas(namespace).List(metav1.ListOptions{})
}