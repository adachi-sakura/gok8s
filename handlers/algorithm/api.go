package algorithm

import (
	"context"
	"errors"
	"fmt"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/model"
	"github.com/buzaiguna/gok8s/prom"
	"github.com/buzaiguna/gok8s/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"math"
	"net/http"
	"time"
)

const (
	nodeMasterRole	= "node-role.kubernetes.io/master"
	dependencyLabel = "dependencies"
	leastResponseTimeLabel	= "leastResponseTime"
	allocationDeploymentLabel = "container-allocation-deployment"
	httpRequestCountLabel = "httpRequestCount"
	tenMinsDuration	= "10m"
)

func GetParameters(ctx context.Context) error {
	newCtx, params, err := buildAlgorithmParameters(ctx)
	if err != nil {
		return err
	}
	appctx.JSON(newCtx, http.StatusOK, params)
	return nil
}

func CreateDeployments(ctx context.Context) error {
	newCtx, params, err := buildAlgorithmParameters(ctx)
	if err != nil {
		return err
	}
	ctx = newCtx

	algorithmCli := appctx.NewAlgorithmClient(ctx)
	allocations, err := algorithmCli.GetAllocations(params)
	if err != nil {
		return err
	}
	deployments := newDeploymentsFromAllocations(ctx, allocations)
	k8sCli := appctx.NewK8SClient(ctx)
	createdDeployments, err := k8sCli.CreateDeployments(deployments)
	if err != nil && len(createdDeployments) == 0 {
		return err
	}
	appctx.JSON(ctx, http.StatusCreated, createdDeployments)
	return nil
}

func newDeploymentsFromAllocations(ctx context.Context, allocations []model.MicroserviceAllocation) []*appsv1.Deployment {
	newDeployments := []*appsv1.Deployment{}
	deployments := appctx.Deployments(ctx)
	for _, microservice := range allocations {
		index, _ := appctx.GetDeploymentIndex(ctx, microservice.Name)
		deployment := deployments[index]
		for num, container := range microservice.Containers {
			newDeployment := *deployment
			newDeployment.Name = fmt.Sprintf("%s-%d", deployment.Name, num)
			newDeployment.Spec.Replicas = utils.NewInt32(1)
			newDeployment.Spec.Template.Spec.NodeSelector = map[string]string{}
			newDeployment.Spec.Template.Spec.NodeName = container.Loc
			newDeployment.Spec.Selector.MatchLabels[allocationDeploymentLabel] = newDeployment.Name
			newDeployment.Spec.Template.Labels[allocationDeploymentLabel] = newDeployment.Name
			requests := v1.ResourceList{}
			requests[v1.ResourceCPU] = *resource.NewMilliQuantity(int64(math.Ceil(container.Cpu)), resource.DecimalSI)
			requests[v1.ResourceMemory] = *resource.NewQuantity(utils.Int64(microservice.RequestMemory).MBtoB(), resource.BinarySI)
			newDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = requests
			newDeployment.Spec.Template.Spec.Containers[0].Resources.Limits = nil

			newDeployments = append(newDeployments, &newDeployment)
		}
	}
	return newDeployments
}

type buildFunc func(context.Context, *model.AlgorithmParameters) error

func buildAlgorithmParameters(ctx context.Context) (context.Context, *model.AlgorithmParameters, error) {
	ctx = buildCliContext(ctx)
	ctx = buildDeploymentsContext(ctx)
	funcs := buildingPipeline()
	metrics := &model.AlgorithmParameters{}
	for _, fun := range funcs {
		if err := fun(ctx, metrics); err != nil {
			return nil, nil, err
		}
	}
	return ctx, metrics, nil
}

func buildCliContext(ctx context.Context) context.Context {
	newCtx := appctx.K8SClientContext(ctx)
	newCtx = appctx.MetricsClientContext(newCtx)
	newCtx = appctx.PromClientContext(newCtx)
	return newCtx
}

func buildDeploymentsContext(ctx context.Context) context.Context {
	ctx = appctx.MultiK8SResourceContext(ctx)
	deployments := appctx.DeploymentObjects(ctx)
	ctx = appctx.WithDeployments(ctx, deployments)
	ctx = appctx.DeploymentInvertedIndexContext(ctx, deployments)
	return ctx
}

func buildTotalTimeRequired(ctx context.Context, metrics *model.AlgorithmParameters) error {
	totalTimeRequired := appctx.Query(ctx, "totalTime")
	if err := metrics.SetTotalTimeRequired(totalTimeRequired); err != nil {
		return apperror.NewInvalidParameterError("totalTime", err)
	}
	return nil
}

func buildLimitRange(ctx context.Context, metrics *model.AlgorithmParameters) error {
	namespace := appctx.Query(ctx, "namespace")
	k8sCli := appctx.K8SClient(ctx)
	limitRanges, err := k8sCli.ListLimitRange(namespace)
	if err != nil {
		return err
	}
	metrics.LimitRange = model.NewLimitRange(limitRanges)
	return nil
}

func buildResourceQuota(ctx context.Context, metrics *model.AlgorithmParameters) error {
	namespace := appctx.Query(ctx, "namespace")
	k8sCli := appctx.NewK8SClient(ctx)
	resourceQuotas, err := k8sCli.ListResourceQuota(namespace)
	if err != nil {
		return err
	}
	metrics.ResourceQuota = model.NewResourceQuota(resourceQuotas)
	return nil
}

func buildNodeParam(ctx context.Context, metrics *model.AlgorithmParameters) error {
	nodesMetrics, err := getNodesMetrics(ctx)
	if err != nil {
		return err
	}
	nodesStatus, err := getNodesStatus(ctx)
	if err != nil {
		return err
	}

	nodesMap := map[string]*model.Node{}
	loadNodesStatus(&nodesMap, nodesStatus)
	loadNodesMetrics(&nodesMap, nodesMetrics)
	nodes := []*model.Node{}
	for _, node := range nodesMap {
		nodes = append(nodes, node)
	}
	metrics.Nodes = nodes

	return nil
}

func loadNodesMetrics(nodesMap *map[string]*model.Node, metricses []v1beta1.NodeMetrics) {
	for _, metrics := range metricses {
		if _, exist := (*nodesMap)[metrics.Name]; !exist {
			continue
		}
		(*nodesMap)[metrics.Name].Current_cpu = metrics.Usage.Cpu().MilliValue()
		(*nodesMap)[metrics.Name].Current_mem = utils.Int64(metrics.Usage.Memory().Value()).BtoMB()
	}
}

func loadNodesStatus(nodesMap *map[string]*model.Node, nodesStatus []v1.Node) {
	for _, nodeStatus := range nodesStatus {
		if _, exist := nodeStatus.Labels[nodeMasterRole]; exist {
			continue
		}
		(*nodesMap)[nodeStatus.Name] = &model.Node{
			Name:	nodeStatus.Name,
			Sum_cpu:	nodeStatus.Status.Capacity.Cpu().MilliValue(),
			Allocatable_cpu:	nodeStatus.Status.Allocatable.Cpu().MilliValue(),
			Sum_mem:	utils.Int64(nodeStatus.Status.Capacity.Memory().Value()).BtoMB(),
			Allocatable_mem:	utils.Int64(nodeStatus.Status.Allocatable.Memory().Value()).BtoMB(),
		}
	}
}

func getNodesMetrics(ctx context.Context) ([]v1beta1.NodeMetrics, error) {
	metricsCli := appctx.MetricsClient(ctx)
	metricsList, err := metricsCli.ListNodeMetrics()
	if err != nil {
		return nil, err
	}
	return metricsList.Items, nil
}

func getNodesStatus(ctx context.Context) ([]v1.Node, error) {
	k8sCli := appctx.K8SClient(ctx)
	nodeList, err := k8sCli.ListNodes()
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

func buildMicroserviceParam(ctx context.Context, metrics *model.AlgorithmParameters) error {
	deployments := appctx.Deployments(ctx)

	if err := validate(ctx); err != nil {
		return err
	}
	datas := []*model.MicroservcieData{}
	for _, deployment := range deployments {
		yamlData, err := buildMicroserviceYaml(ctx, deployment)
		if err != nil {
			return err
		}
		metricsData, err := buildMicroserviceMetrics(ctx, deployment)
		if err != nil {
			return err
		}
		data := &model.MicroservcieData{
			MicroserviceYaml:		*yamlData,
			MicroserviceMetrics:	*metricsData,
		}
		datas = append(datas, data)
	}
	metrics.Datas = datas
	return nil
}

func buildMicroserviceMetrics(ctx context.Context, deployment *appsv1.Deployment) (*model.MicroserviceMetrics, error) {
	metrics := &model.MicroserviceMetrics{}
	t := time.Now()

	funcs := buildingMetricsPipeline()
	for _, fun := range funcs {
		if err := fun(ctx, metrics, deployment, t); err != nil {
			return nil, err
		}
	}
	return metrics, nil
}

type promBuildFunc func(context.Context, *model.MicroserviceMetrics, *appsv1.Deployment, time.Time) error

func buildingMetricsPipeline() []promBuildFunc {
	return []promBuildFunc{
		loadContainerReceiveTotal,
		loadContainerTransmitTotal,
		loadContainerCpuUsageSecTotal,
		//loadHttpRequestsTotal,
		loadContainerMemUsageMax,
	}
}

func loadContainerReceiveTotal(ctx context.Context, metrics *model.MicroserviceMetrics, deployment *appsv1.Deployment, t time.Time) error {
	cli := appctx.PromClient(ctx)
	deployName := deployment.Name
	duration := tenMinsDuration
	res, err := cli.ContainerReceiveTotal(deployName, t, duration)
	if err != nil {
		return err
	}
	matVals := prom.GetMatrixValues(res)
	receive := prom.SumIncrement(matVals...)/1024

	metrics.Receive = receive
	return nil
}

func loadContainerTransmitTotal(ctx context.Context, metrics *model.MicroserviceMetrics, deployment *appsv1.Deployment, t time.Time) error {
	cli := appctx.PromClient(ctx)
	deployName := deployment.Name
	duration := tenMinsDuration
	res, err := cli.ContainerTransmitTotal(deployName, t, duration)
	if err != nil {
		return err
	}
	matVals := prom.GetMatrixValues(res)
	transmit := prom.SumIncrement(matVals...)/1024

	metrics.Transmit = transmit
	return nil
}

func loadContainerCpuUsageSecTotal(ctx context.Context, metrics *model.MicroserviceMetrics, deployment *appsv1.Deployment, t time.Time) error {
	cli := appctx.PromClient(ctx)
	deployName := deployment.Name
	duration := tenMinsDuration
	containerName := deployment.Spec.Template.Spec.Containers[0].Name

	res, err := cli.ContainerCpuUsageSecTotal(deployName, containerName, t, duration)
	if err != nil {
		return err
	}
	matValues := prom.GetMatrixValues(res)
	cpuUsageTime := prom.SumIncrement(matValues...)
	cpuTimeTotal := prom.SumElapsedTime(matValues...)

	metrics.CpuUsageTime = cpuUsageTime
	metrics.CpuTimeTotal = cpuTimeTotal
	return nil
}

//func loadHttpRequestsTotal(ctx context.Context, metrics *model.MicroserviceMetrics, deployment *appsv1.Deployment, t time.Time) error {
//	cli := appctx.PromClient(ctx)
//	deployName := deployment.Name
//	duration := tenMinsDuration
//
//	res, err := cli.HttpRequestsTotal(deployName, t, duration)
//	if err != nil {
//		return err
//	}
//	matValues := prom.GetMatrixValues(res)
//	httpRequestsCount := int(prom.SumIncrement(matValues...))
//
//	metrics.HttpRequestsCount = httpRequestsCount
//	return nil
//}

func loadContainerMemUsageMax(ctx context.Context, metrics *model.MicroserviceMetrics, deployment *appsv1.Deployment, t time.Time) error {
	cli := appctx.PromClient(ctx)
	deployName := deployment.Name
	containerName := deployment.Spec.Template.Spec.Containers[0].Name

	res, err := cli.ContainerMemUsageMax(deployName, containerName, t)
	if err != nil {
		return err
	}
	vecValues := prom.GetVectorValues(res)
	max := prom.Max(vecValues...)
	maxMemoryUsage := float64(max)/1024/1024

	metrics.MaxMemoryUsage = maxMemoryUsage
	return nil
}

func buildMicroserviceYaml(ctx context.Context, deployment *appsv1.Deployment) (*model.MicroserviceYaml, error) {
	yamlData := model.NewMicroserviceYaml()
	yamlData.Name = deployment.Name
	yamlData.Replicas = *deployment.Spec.Replicas
	leastResponseTime := deployment.Annotations[leastResponseTimeLabel]
	if err := yamlData.SetLeastResponseTime(leastResponseTime); err != nil {
		return nil, err
	}
	httpRequestCount := deployment.Annotations[httpRequestCountLabel]
	if err := yamlData.SetHttpRequestCount(httpRequestCount); err != nil {
		return nil, err
	}
	dependencies := deployment.Annotations[dependencyLabel]
	if dependencies != "" {
		dependencyArr := utils.Split(dependencies, ",")
		indexes, err := appctx.GetDeploymentsIndexes(ctx, dependencyArr...)
		if err != nil {
			return nil, err
		}
		yamlData.MicroservicesToInvoke = append(yamlData.MicroservicesToInvoke, indexes...)
	}
	return yamlData, nil
}

func validate(ctx context.Context) error {
	dependenciesMap := map[string][]string{}
	deployments := appctx.Deployments(ctx)
	dict := map[string]bool{}
	//no duplicate deployment
	for _, deployment := range deployments {
		dependencies := deployment.Annotations[dependencyLabel]
		deployName := deployment.Name
		if _, exists := dict[deployName]; exists {
			return apperror.NewInvalidRequestBodyError(errors.New("duplicate deployment name "+deployName))
		}
		dict[deployName] = true
		if dependencies != "" {
			dependencyArr := utils.Split(dependencies, ",")
			dependenciesMap[deployName] = append(dependenciesMap[deployName], dependencyArr...)
		}
	}

	//depended microservice exists (no loop dependency check)
	for _, dependencies := range dependenciesMap {
		for _, nextDeploy := range dependencies {
			_, exists := dict[nextDeploy]
			if !exists {
				return apperror.NewInvalidRequestBodyError(errors.New("depended deploy not found "+nextDeploy))
			}
		}
	}

	return nil
}

func buildBandwidth(ctx context.Context, metrics *model.AlgorithmParameters) error {
	metrics.Bandwidth = config.Bandwidth
	return nil
}

func buildEntrance(ctx context.Context, metrics *model.AlgorithmParameters) error {
	entrance := appctx.Query(ctx, "entry")
	if entrance == "" {
		return apperror.NewParameterRequiredError("entry query")
	}
	_, err := appctx.GetDeploymentIndex(ctx, entrance)
	if err != nil {
		return err
	}
	return nil
}

func buildingPipeline() []buildFunc {
	return []buildFunc{
		buildTotalTimeRequired,
		buildLimitRange,
		buildResourceQuota,
		buildNodeParam,
		buildMicroserviceParam,
		buildEntrance,
		buildBandwidth,
	}
}