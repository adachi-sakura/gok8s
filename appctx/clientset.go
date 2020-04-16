package appctx

import (
	"context"
	"github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/prom"
)

func NewK8SClient(ctx context.Context) *config.K8SClient {
	cfg := RbacConfig(ctx)
	clientSet := config.NewK8SClient(cfg)
	return clientSet
}

func K8SClientContext(ctx context.Context) context.Context {
	cli := K8SClient(ctx)
	if cli != nil {
		return ctx
	}
	cli = NewK8SClient(ctx)
	newCtx := WithK8SClient(ctx, cli)
	return newCtx
}

func NewMetricsClient(ctx context.Context) *config.MetricsClient {
	cfg := RbacConfig(ctx)
	clientSet := config.NewMetricsClient(cfg)
	return clientSet
}

func MetricsClientContext(ctx context.Context) context.Context {
	cli := MetricsClient(ctx)
	if cli != nil {
		return ctx
	}
	cli = NewMetricsClient(ctx)
	newCtx := WithMetricsClient(ctx, cli)
	return newCtx
}

func PromClientContext(ctx context.Context) context.Context {
	cli := PromClient(ctx)
	if cli != nil {
		return ctx
	}
	newCli := prom.PrometheusClient()
	newCtx := WithPromClient(ctx, &newCli)
	return newCtx
}
