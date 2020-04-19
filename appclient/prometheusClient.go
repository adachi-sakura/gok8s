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

func (this *PromClient) ContainerReceiveTotal(deployName string, t time.Time, duration string) (model.Value, error) {
	query := fmt.Sprintf("container_network_receive_bytes_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerTransmitTotal(deployName string, t time.Time, duration string) (model.Value, error) {
	query := fmt.Sprintf("container_network_transmit_bytes_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerCpuUsageSecTotal(deployName string, containerName string, t time.Time, duration string) (model.Value, error) {
	query := fmt.Sprintf("container_cpu_usage_seconds_total{ container =~ \"%s.*\", pod =~ \"%s.*\"}[%s]",
		containerName, deployName, duration)
	return this.Query(query, t)
}

func (this *PromClient) HttpRequestsTotal(deployName string, t time.Time, duration string) (model.Value, error) {
	query := fmt.Sprintf("http_requests_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
	return this.Query(query, t)
}

func (this *PromClient) ContainerMemUsageMax(deployName string, containerName string,t time.Time) (model.Value, error) {
	query := fmt.Sprintf("container_memory_max_usage_bytes{ container =~ \"%s.*\", pod =~ \"%s.*\"}", containerName, deployName)
	return this.Query(query, t)
}

func (this *PromClient) Query(query string, t time.Time) (model.Value, error) {
	res, _, err := v1.NewAPI(this.client).Query(context.Background(), query, t)
	return res, err
}