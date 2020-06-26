package GAMSP

import (
	"context"
	"fmt"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/model"
	"github.com/buzaiguna/gok8s/utils"
	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"log"
	"math"
	"net/http"
)

const (
	annotationBasePrice = "GAMSP/basePrice"
	annotationUnitPrice = "GAMSP/unitPrice"
	nodeMasterRole	= "node-role.kubernetes.io/master"
	allocationDeploymentLabel = "GAMSP-deployment"
	templateDeploymentYaml = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: busybox
spec:
  replicas: 1
  selector:
    matchLabels:
      app: busybox
  template:
    metadata:
      labels:
        app: busybox
    spec:
      containers:
      - name: app
        image: busybox        #内置的linux大多数命令，多用于测试
        imagePullPolicy: IfNotPresent
        args:
        - /bin/sh
        - -c
        - sleep 10; touch /tmp/healthy; sleep 30000
        readinessProbe:           #就绪探针
          exec:
            command:
            - cat
            - /tmp/healthy
          initialDelaySeconds: 10         #10s之后开始第一次探测
          periodSeconds: 5                #第一次探测之后每隔5s探测一次
`
)

var templateDeployment *appsv1.Deployment

func init() {
	obj := utils.ParseK8SYaml([]byte(templateDeploymentYaml))[0]
	templateDeployment = obj.(*appsv1.Deployment)
}

func buildMicroserviceDeployments(allocations []model.GamspAllocation) []*appsv1.Deployment {
	deployments := []*appsv1.Deployment{}
	for msNum, allocation := range allocations {
		for instanceNum, pod := range allocation.Pods {
			newDeployment := *templateDeployment.DeepCopy()
			newDeployment.Name = fmt.Sprintf("%s-%d-%d", "ms", msNum, instanceNum)
			newDeployment.Spec.Replicas = utils.NewInt32(1)
			newDeployment.Spec.Template.Spec.NodeSelector = map[string]string{}
			newDeployment.Spec.Template.Spec.NodeName = pod.Loc
			newDeployment.Spec.Selector.MatchLabels[allocationDeploymentLabel] = newDeployment.Name
			newDeployment.Spec.Template.Labels[allocationDeploymentLabel] = newDeployment.Name
			requests := v1.ResourceList{}
			requests[v1.ResourceCPU] = *resource.NewMilliQuantity(int64(math.Ceil(pod.Cpu)), resource.DecimalSI)
			newDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = requests

			deployments = append(deployments, &newDeployment)

		}
	}
	return deployments
}

func AnnotateNodes(ctx context.Context) error {
	costs := []model.Cost{}
	if err := appctx.Bind(ctx, &costs); err != nil {
		return err
	}
	cli := appctx.NewK8SClient(ctx)
	for _, cost := range costs {
		nodeInfo, err := cli.GetNode(cost.Name)
		if err != nil {
			return err
		}
		nodeInfo.Annotations[annotationBasePrice] = cost.BasePrice
		nodeInfo.Annotations[annotationUnitPrice] = cost.UnitPrice
		if _, err := cli.UpdateNode(nodeInfo); err != nil {
			return err
		}
	}
	appctx.JSON(ctx, http.StatusOK, nil)
	return nil
}

func CreateDeployments(ctx context.Context) error {
	k8sCli := appctx.NewK8SClient(ctx)
	nodesList, err := k8sCli.ListNodes()
	if err != nil {
		return err
	}
	nodesInfo := nodesList.Items
	nodes := []*model.GamspNodeInfo{}
	ctx = appctx.MetricsClientContext(ctx)
	for _, nodeInfo := range nodesInfo {
		if _, exist := nodeInfo.Labels[nodeMasterRole]; exist {
			continue
		}
		gamspNode, err := buildGamspNode(ctx, nodeInfo)
		if err != nil {
			return err
		}
		nodes = append(nodes, gamspNode)
	}
	params := &model.GamspParameters{
		Nodes:	nodes,
	}
	log.Println(params)

	algorithmCli := appctx.NewAlgorithmClient(ctx)
	allocations, err := algorithmCli.GetGamspAllocations(params)
	if err != nil {
		return err
	}
	log.Printf("Optimized container allocation result:\n%v\n",allocations)
	deployments := buildMicroserviceDeployments(allocations)
	createdDeployments, err := k8sCli.CreateDeployments(deployments)
	if err != nil && len(createdDeployments) == 0 {
		return err
	}
	appctx.JSON(ctx, http.StatusCreated, createdDeployments)
	return nil
}

func buildGamspNode(ctx context.Context, nodeInfo v1.Node) (*model.GamspNodeInfo, error) {
	if _, exist := nodeInfo.Annotations[annotationBasePrice]; !exist {
		return nil, apperror.NewInternalServerError(fmt.Sprintf("annotation %s not found", annotationBasePrice))
	}
	if _, exist := nodeInfo.Annotations[annotationUnitPrice]; !exist {
		return nil, apperror.NewInternalServerError(fmt.Sprintf("annotation %s not found", annotationUnitPrice))
	}
	nodeMetric, err := getNodeMetrics(ctx, nodeInfo.Name)
	if err != nil {
		return nil, err
	}
	return &model.GamspNodeInfo{
		Name:      nodeInfo.Name,
		MilliCore: nodeInfo.Status.Allocatable.Cpu().MilliValue()-nodeMetric.Usage.Cpu().MilliValue(),
		Mem:	   utils.Int64(nodeInfo.Status.Allocatable.Memory().Value()-nodeMetric.Usage.Memory().Value()).BtoMB(),
		BasePrice: utils.MustParseFloat64(nodeInfo.Annotations[annotationBasePrice]),
		UnitPrice: utils.MustParseFloat64(nodeInfo.Annotations[annotationUnitPrice]),
	}, nil
}

func getNodeMetrics(ctx context.Context, name string) (*v1beta1.NodeMetrics, error) {
	metricsCli := appctx.MetricsClient(ctx)
	return metricsCli.GetNodeMetrics(name)
}