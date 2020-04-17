package appctx

import (
	"context"
	"github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/model"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

const (
	keyGinContext 	= "GinContext"
	keyContext		= "ServiceContext"
	keyInClusterConfig = "inClusterConfig"
	keyK8SToken		= "k8sToken"
	keyK8SClient	= "k8sClient"
	keyRBACConfig	= "rbacConfig"
	keyK8SObjects	= "k8sObjects"
	keyMetricsClient = "metricsClient"
	keyPromClient	= "promClient"
	keyNodesMap		= "nodesMap"
	keyDeployments	= "deployments"
	keyDeploymentIndex = "deploymentIndex"
)

func WithGinContext(ctx context.Context, ginContext *gin.Context) context.Context {
	return context.WithValue(ctx, keyGinContext, ginContext)
}

func JSON(ctx context.Context, code int, obj interface{}) {
	c := GinContext(ctx)
	c.JSON(code, obj)
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

func WithRbacConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, keyRBACConfig, cfg)
}

func RbacConfig(ctx context.Context) *rest.Config {
	val := ctx.Value(keyRBACConfig)
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

func WithK8SClient(ctx context.Context, client *config.K8SClient) context.Context {
	return context.WithValue(ctx, keyK8SClient, client)
}

func K8SClient(ctx context.Context) *config.K8SClient {
	val := ctx.Value(keyK8SClient)
	if val == nil {
		return nil
	}
	return val.(*config.K8SClient)
}

func WithMetricsClient(ctx context.Context, client *config.MetricsClient) context.Context {
	return context.WithValue(ctx, keyMetricsClient, client)
}

func MetricsClient(ctx context.Context) *config.MetricsClient {
	val := ctx.Value(keyMetricsClient)
	if val == nil {
		return nil
	}
	return val.(*config.MetricsClient)
}

func WithPromClient(ctx context.Context, client *config.PromClient) context.Context {
	return context.WithValue(ctx, keyPromClient, client)
}

func PromClient(ctx context.Context) *config.PromClient{
	val := ctx.Value(keyPromClient)
	if val == nil {
		return nil
	}
	return val.(*config.PromClient)
}

func WithK8SObjects(ctx context.Context, objects []runtime.Object) context.Context {
	return context.WithValue(ctx, keyK8SObjects, objects)
}

func K8SObjects(ctx context.Context) []runtime.Object {
	val := ctx.Value(keyK8SObjects)
	if val == nil {
		return nil
	}
	return val.([]runtime.Object)
}

func WithNodesMap(ctx context.Context, m map[string]*model.Node) context.Context {
	return context.WithValue(ctx, keyNodesMap, m)
}

func NodesMap(ctx context.Context) map[string]*model.Node {
	val := ctx.Value(keyNodesMap)
	if val == nil {
		return nil
	}
	return val.(map[string]*model.Node)
}

func WithDeployments(ctx context.Context, deployments []*v1.Deployment) context.Context {
	return context.WithValue(ctx, keyDeployments, deployments)
}

func Deployments(ctx context.Context) []*v1.Deployment {
	val := ctx.Value(keyDeployments)
	if val == nil {
		return nil
	}
	return val.([]*v1.Deployment)
}

func WithDeploymentIndex(ctx context.Context, mNameToIndex map[string]int) context.Context {
	return context.WithValue(ctx, keyDeploymentIndex, mNameToIndex)
}

func DeploymentIndex(ctx context.Context) map[string]int {
	val := ctx.Value(keyDeploymentIndex)
	if val == nil {
		return nil
	}
	return val.(map[string]int)
}

