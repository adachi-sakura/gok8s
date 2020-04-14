package auth

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/config"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
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
	cli := generateClient(ctx)
	newCtx := appctx.WithK8SClient(ctx, cli)
	appctx.SetContext(c, newCtx)
	return nil
}

func generateClient(ctx context.Context) *kubernetes.Clientset {
	cfg := appctx.InClusterConfig(ctx)
	token := appctx.K8SToken(ctx)
	cfg.BearerTokenFile = ""
	cfg.BearerToken = token
	clientSet := config.NewClient(cfg)
	return clientSet
}


var AuthHandler = appctx.GinHandler(Authentication)