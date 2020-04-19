package appclient

import (
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

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