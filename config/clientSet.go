package config

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (cli *K8SClient) ListPods(namespace string) (*v1.PodList, error){
	return cli.K8SClient().CoreV1().Pods(namespace).List(metav1.ListOptions{})
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

func (cli *MetricsClient) GetNodeMetrics(name string) (*v1beta1.NodeMetrics, error) {
	return cli.MetricsClient().MetricsV1beta1().NodeMetricses().Get(name, metav1.GetOptions{})
}