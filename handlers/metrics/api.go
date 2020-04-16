package metrics

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"net/http"
)

func GetNodesMetrics(ctx context.Context) error {
	cli := appctx.NewMetricsClient(ctx)
	metricses, err := cli.ListNodeMetrics()
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusOK, metricses)
	return nil
}

func GetNodeMetrics(ctx context.Context) error {
	nodeName := appctx.Param(ctx, "nodeName")
	cli := appctx.NewMetricsClient(ctx)
	metrics, err := cli.GetNodeMetrics(nodeName)
	if err != nil {
		return apperror.NewInvalidParameterError(nodeName, err)
	}
	appctx.JSON(ctx, http.StatusOK, metrics)
	return nil
}

func GetPodsMetrics(ctx context.Context) error {
	cli := appctx.NewMetricsClient(ctx)
	namespace := appctx.Query(ctx, "namespace")
	metricses, err := cli.ListPodMetrics(namespace)
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusOK, metricses)
	return nil
}
