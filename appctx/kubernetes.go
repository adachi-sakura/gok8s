package appctx

import (
	"context"
	"encoding/json"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/utils"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	NULL	= 0
)

func BindK8SYaml(ctx context.Context, obj interface{}) error {
	c := GinContext(ctx)
	jsonBody := getYamlBody(c).yamlToJson()
	if err := json.Unmarshal(jsonBody, obj); err != nil {
		return apperror.NewInvalidRequestBodyError(err)
	}
	return nil
}

func MultiK8SResourceContext(ctx context.Context) context.Context {
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

func DeploymentInvertedIndexContext(ctx context.Context, deployments []*v1.Deployment) context.Context {
	mNameToIndex := map[string]int{}
	for num, deployment := range deployments {
		mNameToIndex[deployment.Name] = num
	}
	newCtx := WithDeploymentIndex(ctx, mNameToIndex)
	return newCtx
}

func GetDeploymentIndex(ctx context.Context, deploymentName string) (int, error) {
	mNameToIndex := DeploymentIndex(ctx)
	if mNameToIndex == nil {
		return NULL, apperror.NewInternalServerError("deployment map is nil")
	}
	index, exists := mNameToIndex[deploymentName]
	if !exists {
		return NULL, apperror.NewResourceNotFoundError("deployment "+deploymentName)
	}
	return index, nil
}

func GetDeploymentsIndexes(ctx context.Context, deploymentNames ...string) ([]int, error) {
	indexes := []int{}
	mNameToIndex := DeploymentIndex(ctx)
	if mNameToIndex == nil {
		return nil, apperror.NewInternalServerError("deployment map is nil")
	}
	for _, name := range deploymentNames {
		index, exists := mNameToIndex[name]
		if !exists {
			return nil, apperror.NewResourceNotFoundError("deployment "+name)
		}
		indexes = append(indexes, index)
	}

	return indexes, nil
}