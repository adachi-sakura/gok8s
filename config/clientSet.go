package config

import (
	"github.com/prometheus/client_golang/api"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
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

func (cli *K8SClient) ListPods(namespace string) (*corev1.PodList, error){
	return cli.K8SClient().CoreV1().Pods(namespace).List(metav1.ListOptions{})
}

func (cli *K8SClient) GetNode(name string) (*corev1.Node, error) {
	return cli.K8SClient().CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (cli *K8SClient) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	return cli.K8SClient().AppsV1().Deployments(namespace).Create(deployment)
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

type MetricsClient metrics.Clientset

func NewMetricsClient(config *rest.Config) *MetricsClient {
	clientSet, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return (*MetricsClient)(clientSet)
}

func (cli *MetricsClient) MetricsClient() *metrics.Clientset {
	return (*metrics.Clientset)(cli)
}

func (cli *MetricsClient) ListNodeMetrics() (*v1beta1.NodeMetricsList, error) {
	return cli.MetricsClient().MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
}

func (cli *MetricsClient) ListPodMetrics(namespace string ) (*v1beta1.PodMetricsList, error) {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	return cli.MetricsClient().MetricsV1beta1().PodMetricses(namespace).List(metav1.ListOptions{})
}

func (cli *MetricsClient) GetNodeMetrics(name string) (*v1beta1.NodeMetrics, error) {
	return cli.MetricsClient().MetricsV1beta1().NodeMetricses().Get(name, metav1.GetOptions{})
}
