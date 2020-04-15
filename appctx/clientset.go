package appctx

import (
	"context"
	"github.com/buzaiguna/gok8s/config"
)

func NewK8SClient(ctx context.Context) *config.K8SClient {
	cfg := RbacConfig(ctx)
	clientSet := config.NewK8SClient(cfg)
	return clientSet
}

func NewMetricsClient(ctx context.Context) *config.MetricsClient {
	cfg := RbacConfig(ctx)
	clientSet := config.NewMetricsClient(cfg)
	return clientSet
}
