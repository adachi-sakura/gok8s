package appctx

import (
	"context"
	"encoding/json"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/utils"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

func BindK8SYaml(ctx context.Context, obj interface{}) error {
	c := GinContext(ctx)
	jsonBody := getYamlBody(c).yamlToJson()
	if err := json.Unmarshal(jsonBody, obj); err != nil {
		return apperror.NewInvalidRequestBodyError(err)
	}
	return nil
}

func DecodeMultiK8SResource(ctx context.Context) context.Context {
	c := GinContext(ctx)
	yamlFiles := getYamlBody(c)
	objects := utils.ParseK8SYaml(yamlFiles)
	return WithK8SObjects(ctx, objects)
}

func DeploymentObjects(ctx context.Context) []*v1.Deployment {
	objects := K8SObjects(ctx)
	deployments := []*v1.Deployment{}
	for _, obj := range objects {
		switch o := obj.(type) {
		case *v1.Deployment:
			deployments = append(deployments, o)
		default:
		}
	}
	return deployments
}

func ServiceObjects(ctx context.Context) []*corev1.Service {
	objects := K8SObjects(ctx)
	services := []*corev1.Service{}
	for _, obj := range objects {
		switch o := obj.(type) {
		case *corev1.Service:
			services = append(services, o)
		default:
		}
	}
	return services
}
