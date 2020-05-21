package appclient

import (
	"context"
	"fmt"
	"github.com/buzaiguna/gok8s/prom"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

const (
	containerMetricsQueryFormat = "%s{ container = \"%s\", pod =~ \"%s.*\"}"
	additionalRangeFormat = "[%s]"
)

type PromClient struct {
	client api.Client
}

func NewPromClient() *PromClient {
	return &PromClient{
		client:	prom.PrometheusClient(),
	}
}

func (this *PromClient) PromClient() api.Client {
	return this.client
}

func (this *PromClient) ContainerReceiveTotal(deployName string, containerName string, t time.Time, duration string) (model.Value, error) {
	query := buildContainerMetricsRangeQuery("container_network_receive_bytes_total", deployName, containerName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerTransmitTotal(deployName string, containerName string, t time.Time, duration string) (model.Value, error) {
	query := buildContainerMetricsRangeQuery("container_network_transmit_bytes_total", deployName, containerName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerCpuUsageSecTotal(deployName string, containerName string, t time.Time, duration string) (model.Value, error) {
	query := buildContainerMetricsRangeQuery("container_cpu_usage_seconds_total", deployName, containerName, duration)
	return this.Query(query, t)
}

func (this *PromClient) HttpRequestsTotal(deployName string, containerName string, t time.Time, duration string) (model.Value, error) {
	query := buildContainerMetricsRangeQuery("http_requests_total", deployName, containerName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerMemUsageMax(deployName string, containerName string,t time.Time) (model.Value, error) {
	query := buildContainerMetricsQuery("container_memory_max_usage_bytes", deployName, containerName)
	return this.Query(query, t)
}

func (this *PromClient) Query(query string, t time.Time) (model.Value, error) {
	res, _, err := v1.NewAPI(this.client).Query(context.Background(), query, t)
	return res, err
}

func buildContainerMetricsQuery(metricsName string, deployName string, containerName string) string {
	return fmt.Sprintf(containerMetricsQueryFormat, metricsName, containerName, deployName)
}

func buildContainerMetricsRangeQuery(metricsName string, deployName string, containerName string, duration string) string {
	containerMetricsRangeQueryFormat := containerMetricsQueryFormat+additionalRangeFormat
	return fmt.Sprintf(containerMetricsRangeQueryFormat, metricsName, deployName, containerName, duration)
}