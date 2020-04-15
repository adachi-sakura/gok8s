package auth

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
)

func validateToken(c *gin.Context) error {
	token := c.Request.Header.Get("Token")
	if token == "" {
		return apperror.NewHeaderRequiredError("Token")
	}
	ctx := appctx.GetContextFromGin(c)
	newCtx := appctx.WithK8SToken(ctx, token)
	appctx.SetContext(c, newCtx)
	return nil
}

func Authentication(c *gin.Context)  error {
	if err := validateToken(c); err != nil {
		return err
	}
	ctx := appctx.GetContextFromGin(c)
	rbacConfig := generateRbacConfig(ctx)
	newCtx := appctx.WithRbacConfig(ctx, rbacConfig)
	//cli := appctx.NewK8SClient(newCtx)
	//newCtx = appctx.WithK8SClient(newCtx, cli)
	appctx.SetContext(c, newCtx)
	return nil
}

func generateRbacConfig(ctx context.Context) *rest.Config {
	cfg := appctx.InClusterConfig(ctx)
	token := appctx.K8SToken(ctx)
	cfg.BearerTokenFile = ""
	cfg.BearerToken = token
	return cfg
}

var AuthHandler = appctx.GinHandler(Authentication)