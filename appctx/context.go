package appctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	keyGinContext 	= "GinContext"
	keyContext		= "ServiceContext"
	keyInClusterConfig = "inClusterConfig"
	keyK8SToken		= "k8sToken"
	keyK8SClient	= "k8sClient"
)

func WithGinContext(ctx context.Context, ginContext *gin.Context) context.Context {
	return context.WithValue(ctx, keyGinContext, ginContext)
}

func GinContext(ctx context.Context) *gin.Context {
	val := ctx.Value(keyGinContext)
	if val == nil {
		return nil
	}
	return val.(*gin.Context)
}

func GetContextFromGin(c *gin.Context) context.Context {
	val, _ := c.Get(keyContext)
	return val.(context.Context)
}

func SetContext(ginContext *gin.Context, ctx context.Context) {
	ginContext.Set(keyContext, ctx)
}

func WithInClusterConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, keyInClusterConfig, cfg)
}

func InClusterConfig(ctx context.Context) *rest.Config {
	val := ctx.Value(keyInClusterConfig)
	if val == nil {
		return nil
	}
	return val.(*rest.Config)
}

func WithK8SToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, keyK8SToken, token)
}

func K8SToken(ctx context.Context) string {
	val := ctx.Value(keyK8SToken)
	if val == nil {
		return ""
	}
	return val.(string)
}

func WithK8SClient(ctx context.Context, client *kubernetes.Clientset) context.Context {
	return context.WithValue(ctx, keyK8SClient, client)
}

func K8SClient(ctx context.Context) *kubernetes.Clientset {
	val := ctx.Value(keyK8SClient)
	if val == nil {
		return nil
	}
	return val.(*kubernetes.Clientset)
}