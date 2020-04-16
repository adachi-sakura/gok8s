package algorithm

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/model"
	"net/http"
)

func GetMetrics(ctx context.Context) error {
	newCtx, metrics, err := buildAlgorithmParameters(ctx)
	if err != nil {
		return err
	}
	appctx.JSON(newCtx, http.StatusOK, metrics)
	return nil
}

type buildFunc func(context.Context, *model.AlgorithmParameters) (context.Context, error)

func buildAlgorithmParameters(ctx context.Context) (context.Context, *model.AlgorithmParameters, error) {
	funcs := buildingPipeline()
	metrics := &model.AlgorithmParameters{}
	var err error
	for _, fun := range funcs {
		ctx, err = fun(ctx, metrics)
		if err != nil {
			return nil, nil, err
		}
	}
	return ctx, metrics, nil
}

func buildCli(ctx context.Context, metrics *model.AlgorithmParameters) (context.Context, error) {
	newCtx := appctx.K8SClientContext(ctx)
	newCtx = appctx.MetricsClientContext(newCtx)
	newCtx = appctx.PromClientContext(newCtx)
	return newCtx, nil
}

func buildTotalTimeRequired(ctx context.Context, metrics *model.AlgorithmParameters) (context.Context, error) {
	totalTimeRequired := appctx.Query(ctx, "totalTime")
	if err := metrics.SetTotalTimeRequired(totalTimeRequired); err != nil {
		return nil, apperror.NewInvalidParameterError("totalTime", err)
	}
	return ctx, nil
}

func buildLimitRange(ctx context.Context, metrics *model.AlgorithmParameters) (context.Context, error) {
	namespace := appctx.Query(ctx, "namespace")
	k8sCli := appctx.K8SClient(ctx)
	limitRanges, err := k8sCli.ListLimitRange(namespace)
	if err != nil {
		return nil, err
	}
	metrics.LimitRange = model.NewLimitRange(limitRanges)
	return ctx, nil
}

func buildResourceQuota(ctx context.Context, metrics *model.AlgorithmParameters) (context.Context, error) {
	namespace := appctx.Query(ctx, "namespace")
	k8sCli := appctx.NewK8SClient(ctx)
	resourceQuotas, err := k8sCli.ListResourceQuota(namespace)
	if err != nil {
		return nil, err
	}
	metrics.ResourceQuota = model.NewResourceQuota(resourceQuotas)
	return ctx, nil
}

func buildingPipeline() []buildFunc {
	return []buildFunc{
		buildCli,
		buildTotalTimeRequired,
		buildLimitRange,
		buildResourceQuota,
	}
}