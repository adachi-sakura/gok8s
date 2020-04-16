package apiserver

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	v1 "k8s.io/api/apps/v1"
	"net/http"
)

func GetPods(ctx context.Context) error {
	cli := appctx.NewK8SClient(ctx)
	namespace := appctx.Query(ctx, "namespace")
	pods, err := cli.ListPods(namespace)
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusOK, pods)
	return nil
}

func GetNode(ctx context.Context) error {
	nodeName := appctx.Param(ctx, "nodeName")
	cli := appctx.NewK8SClient(ctx)
	node, err := cli.GetNode(nodeName)
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusOK, node)
	return nil
}

func CreateDeployment(ctx context.Context) error {
	deployment := &v1.Deployment{}
	if err := appctx.BindK8SYaml(ctx, deployment); err != nil {
		return err
	}
	cli := appctx.NewK8SClient(ctx)
	result, err := cli.CreateDeployment(deployment.Namespace, deployment)
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusCreated, result)
	return nil
}

