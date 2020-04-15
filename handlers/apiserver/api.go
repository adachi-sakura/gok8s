package apiserver

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"net/http"
)

func GetPods(ctx context.Context) error {
	cli := appctx.K8SClient(ctx)
	namespace := appctx.Param(ctx, "namespace")
	pods, err := cli.ListPods(namespace)
	if err != nil {
		return err
	}
	appctx.JSON(ctx, http.StatusOK, pods)
	return nil
}


