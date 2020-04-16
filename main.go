package main
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cfg "github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/handlers"
	"github.com/buzaiguna/gok8s/middleware"
	"github.com/buzaiguna/gok8s/middleware/basic"
	"github.com/buzaiguna/gok8s/model"
	myProm "github.com/buzaiguna/gok8s/prom"
	"github.com/buzaiguna/gok8s/utils"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultCmdConfigName = "kubernetes"
)

func main() {
	listenAndServe()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("initial config token: %s\ntokenfile: %s", config.BearerToken, config.BearerTokenFile)
	clientSet := &kubernetes.Clientset{}

	//todo unit test
	type mapItem struct {
		deployments map[string]*appsv1.Deployment
		nodes	map[string]*model.Node
	}
	deploymentsMap := map[string]mapItem{}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "PONG")
	})
	router.GET("/admin/first-pod", func(c *gin.Context) {
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println("error occurred in cluster config...")
			panic(err)
		}
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println("error occurred in clientSet creating...")
			panic(err)
		}
		fmt.Printf("ClientSet: %v\n", clientSet)
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, pods.Items[0])

	})
	router.GET("/pods", DynamicClientSet(config, &clientSet), func(c *gin.Context) {
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pods)

	})
	router.GET("/first-pod", DynamicClientSet(config, &clientSet), func(c *gin.Context) {
		fmt.Printf("clientSet used: %v\n", *clientSet)
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pods.Items[0])

	})
	router.GET("/test/first-pod", DyClientSet(config, &clientSet), func(c *gin.Context) {
		fmt.Printf("clientSet used: %v\n", *clientSet)
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pods.Items[0])

	})
	router.GET("/metrics/nodes", DyClientSet(config, &clientSet), func(c *gin.Context) {
		metricsClientSet, err := metricsclientset.NewForConfig(config)
		if err != nil {
			fmt.Printf("error occurred in metrics client set creating...")
			panic(err)
		}
		metrics, err := metricsClientSet.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("error occurred in metrics list...")
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, metrics)

	})
	router.GET("/metrics/nodes/:nodeName", DyClientSet(config, &clientSet), func(c *gin.Context) {
		metricsClientSet, err := metricsclientset.NewForConfig(config)
		if err != nil {
			fmt.Printf("error occurred in metrics client set creating...")
			panic(err)
		}
		name := c.Param("nodeName")
		metrics, err := metricsClientSet.MetricsV1beta1().NodeMetricses().Get(name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("error occurred in metrics get...")
			c.JSON(http.StatusBadRequest, err.Error())
		}
		fmt.Println(metrics.Usage[apiv1.ResourceCPU])
		c.JSON(http.StatusOK, metrics)

	})
	router.GET("/nodes/:nodeName", DyClientSet(config, &clientSet), func(c *gin.Context) {

		name := c.Param("nodeName")
		node, _ := clientSet.CoreV1().Nodes().Get(name, metav1.GetOptions{})
		fmt.Println(node.Status.Capacity[apiv1.ResourceCPU])
		fmt.Println(node.Status.Allocatable[apiv1.ResourceMemory])
		c.JSON(http.StatusOK, node)

	})
	router.GET("/metrics/pods", DyClientSet(config, &clientSet), func(c *gin.Context) {
		metricsClientSet, err := metricsclientset.NewForConfig(config)
		if err != nil {
			fmt.Printf("error occurred in metrics client set creating...")
			panic(err)
		}
		metrics, err := metricsClientSet.MetricsV1beta1().PodMetricses(metav1.NamespaceDefault).List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("error occurred in metrics list...")
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, metrics)

	})
	router.GET("/prom/api-request-total", func(c *gin.Context) {
		cli := myProm.PrometheusClient()
		t := time.Now()
		r := prom.Range{
			Start:	t.Add(-3*time.Hour),
			End:	t,
			Step:	time.Hour,

		}
		res, _, err := prom.NewAPI(cli).QueryRange(context.Background(), "apiserver_request_total", r)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(200, res)
	})
	router.POST("/deployments", DynamicClientSet(config, &clientSet), func(c * gin.Context) {
		deployment := &appsv1.Deployment{}
		if err := utils.Bind(c, deployment); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		fmt.Printf("deployment received:\n%v\n", deployment)
		fmt.Println("annotations tagged as 'test' is:")
		fmt.Println(deployment.Spec.Template.Annotations["test"])
		result, err := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Create(deployment)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusCreated, result)
	})
	router.POST("/decoding", func(c *gin.Context) {
		services := []*apiv1.Service{}
		deployments := []*appsv1.Deployment{}
		objs, err := utils.DecodeK8SResources(c)
		if err != nil {
			fmt.Println("decoding failed")
			panic(err)
		}
		for _, obj := range objs {
			switch o := obj.(type) {
			case *appsv1.Deployment:
				deployments = append(deployments, o)
			case *apiv1.Service:
				services = append(services, o)
			default:
			}
		}
		ret := []interface{}{}
		for _, service := range services {
			ret = append(ret, service)
		}
		for _, deployment := range deployments {
			ret = append(ret, deployment)
		}
		c.JSON(200, ret)
	})
	router.GET("/all-metrics", DyClientSet(config, &clientSet), func(c *gin.Context) {
		item := mapItem{
			deployments: map[string]*appsv1.Deployment{},
			nodes: map[string]*model.Node{},
		}
		ret := model.AlgorithmParameters{
			TotalTimeRequired:	float64(utils.INT_MAX),
		}
		totalTimeRequired := c.Query("totalTime")
		if totalTimeRequired != "" {
			ret.TotalTimeRequired, err = utils.Float64(totalTimeRequired)
			if err != nil {
				panic(err)
			}
		}
		deployments := []*appsv1.Deployment{}
		objs, err := utils.DecodeK8SResources(c)
		if err != nil {
			fmt.Println("decoding failed")
			panic(err)
		}
		for _, obj := range objs {
			switch o := obj.(type) {
			case *appsv1.Deployment:
				deployments = append(deployments, o)
			default:
			}
		}

		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = apiv1.NamespaceDefault
		}
		limitRanges, err := clientSet.CoreV1().LimitRanges(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		lm := model.LimitRange{int64(utils.INT_MAX), int64(utils.INT_MAX)}
		for _, limitRange := range limitRanges.Items {
			for _, item := range limitRange.Spec.Limits {
				if item.Type != apiv1.LimitTypeContainer {
					continue
				}
				if item.Max == nil {
					continue
				}
				if maxCpu, exists := item.Max[apiv1.ResourceCPU]; exists {
					lm.Cpu_lm = utils.Int64Min(lm.Cpu_lm, maxCpu.MilliValue())
				}
				if maxMem, exists := item.Max[apiv1.ResourceMemory]; exists {
					lm.Mem_lm = utils.Int64Min(lm.Mem_lm, maxMem.Value())
				}

			}
		}
		ret.LimitRange = lm

		resourceQuotas, err := clientSet.CoreV1().ResourceQuotas(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		rq := model.ResourceQuota{int64(utils.INT_MAX), int64(utils.INT_MAX)}
		for _, resourceQuota := range resourceQuotas.Items {
			if resourceQuota.Spec.Hard == nil {
				continue
			}
			if maxCpu, exists := resourceQuota.Spec.Hard[apiv1.ResourceCPU]; exists {
				rq.Cpu_rq_total = utils.Int64Min(rq.Cpu_rq_total, maxCpu.MilliValue())
			}
			if maxMem, exists := resourceQuota.Spec.Hard[apiv1.ResourceMemory]; exists {
				rq.Mem_rq_total = utils.Int64Min(rq.Mem_rq_total, maxMem.Value())
			}
		}
		ret.ResourceQuota = rq

		metricsClientSet, err := metricsclientset.NewForConfig(config)
		if err != nil {
			fmt.Printf("error occurred in metrics client set creating...")
			panic(err)
		}
		metricsList, err := metricsClientSet.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("error occurred in metrics list...")
			c.JSON(http.StatusBadRequest, err.Error())
		}
		nodeMetrics := metricsList.Items
		nodesStatus := []apiv1.Node{}
		nodesList, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		nodesStatus = append(nodesStatus, nodesList.Items...)
		nodes := map[string]*model.Node{}
		for _, nodeStatus := range nodesStatus {
			if _, exist := nodeStatus.Labels["node-role.kubernetes.io/master"]; exist {
				continue
			}
			nodes[nodeStatus.Name] = &model.Node{
				Name:		nodeStatus.Name,
				Sum_cpu:	nodeStatus.Status.Capacity.Cpu().MilliValue(),
				Allocatable_cpu:	nodeStatus.Status.Allocatable.Cpu().MilliValue(),
				Sum_mem:	nodeStatus.Status.Capacity.Memory().Value()/1024/1024,
				Allocatable_mem:	nodeStatus.Status.Allocatable.Memory().Value()/1024/1024,
			}
		}
		for _, metric := range nodeMetrics {
			if _, exist := nodes[metric.Name]; !exist {
				continue
			}
			//fmt.Println("by CPU(): ",metric.Usage.Cpu())
			//fmt.Println("mili: ", metric.Usage.Cpu().MilliValue())
			//fmt.Println("by map: ", metric.Usage[apiv1.ResourceCPU])
			nodes[metric.Name].Current_cpu = metric.Usage.Cpu().MilliValue()
			nodes[metric.Name].Current_mem = metric.Usage.Memory().Value()/1024/1024
		}
		for _, node := range nodes {
			ret.Nodes = append(ret.Nodes, node)
		}
		//todo storage
		item.nodes = nodes
		t := time.Now()
		duration := "10m"
		datas := []*model.MicroservcieData{}
		entrance := c.Query("entry")
		if entrance == "" {
			c.JSON(http.StatusBadRequest, "entrance point required")
			return
		}
		dict := map[string]int{}
		dependencyMap := map[string][]string{}
		for num, deployment := range deployments {
			data := &model.MicroservcieData{
				MicroserviceMetrics: model.MicroserviceMetrics{},
				MicroserviceYaml: model.MicroserviceYaml{
					LeastResponseTime:	float64(utils.INT_MAX),
					MicroservicesToInvoke: []int{},
				},
			}
			data.Replicas = *deployment.Spec.Replicas

			leastResponseTime := deployment.Labels["leastResponseTime"]
			if leastResponseTime != "" {
				data.LeastResponseTime, err = utils.Float64(leastResponseTime)
				if err != nil {
					panic(err.Error())
				}
			} else {
				data.LeastResponseTime = float64(utils.INT_MAX)
			}

			deployName := deployment.Name
			//todo storage
			item.deployments[deployName] = deployment
			data.Name = deployName
			dict[deployName] = num
			dependencies := deployment.Labels["dependencies"]
			if dependencies != "" {
				dependencyArr := strings.Split(dependencies, ",")
				dependencyArr = utils.TrimSpace(dependencyArr)
				dependencyMap[deployName] = append(dependencyMap[deployName], dependencyArr...)
			}

			containerName := deployment.Spec.Template.Spec.Containers[0].Name

			query := fmt.Sprintf("container_network_receive_bytes_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
			res, err := myProm.Query(query, t)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			matValues := myProm.GetMatrixValues(res)
			data.Receive = myProm.SumIncrement(matValues...)/1024

			query = fmt.Sprintf("container_network_transmit_bytes_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
			res, err = myProm.Query(query, t)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			matValues = myProm.GetMatrixValues(res)
			data.Transmit = myProm.SumIncrement(matValues...)/1024

			query = fmt.Sprintf("container_cpu_usage_seconds_total{ container =~ \"%s.*\", pod =~ \"%s.*\"}[%s]",
								containerName, deployName, duration)
			res, err = myProm.Query(query, t)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			matValues = myProm.GetMatrixValues(res)
			data.CpuUsageTime = myProm.SumIncrement(matValues...)
			data.CpuTimeTotal = myProm.SumElapsedTime(matValues...)

			query = fmt.Sprintf("http_requests_total{ pod =~ \"%s.*\"}[%s]", deployName, duration)
			res, err = myProm.Query(query, t)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			matValues = myProm.GetMatrixValues(res)
			data.HttpRequestsCount = int(myProm.SumIncrement(matValues...))

			query = fmt.Sprintf("container_memory_max_usage_bytes{ container =~ \"%s.*\", pod =~ \"%s.*\"}", containerName, deployName)
			res, err = myProm.Query(query, t)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			vecValues := myProm.GetVectorValues(res)
			max := myProm.Max(vecValues...)
			data.MaxMemoryUsage = float64(max)/1024/1024

			//fmt.Println("response type is: "+res.Type().String())
			//fmt.Println("response string is:"+res.String())
			datas = append(datas, data)
		}

		num, exists := dict[entrance]
		if !exists {
			c.JSON(http.StatusBadRequest, "entrance service name not found")
			return
		}
		ret.EntrancePoint = num
		ret.Bandwidth = cfg.Bandwidth


		for serviceName, dependencies := range dependencyMap {
			num := dict[serviceName]
			for _, nextServiceName := range dependencies {
				dependencyNum, exists := dict[nextServiceName]
				if !exists {
					c.JSON(http.StatusBadRequest, "invalid dependency name")
					return
				}
				datas[num].MicroservicesToInvoke = append(datas[num].MicroservicesToInvoke, dependencyNum)
			}
		}

		ret.Datas = datas

		//todo storage
		id := bson.NewObjectId().Hex()
		deploymentsMap[id] = item
		fmt.Println(deploymentsMap)

		fmt.Println(id)

		fmt.Println("Ready to call algorithm")

		jsonBytes, err := json.Marshal(ret)
		if err != nil {
			fmt.Println("marshal failed")
			panic(err)
		}
		url := fmt.Sprintf("http://%s:%s/algorithm", cfg.ALGORITHM_HOST, cfg.ALGORITHM_PORT)
		request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			fmt.Println("create request failed")
			panic(err)
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("request failed")
			panic(err)
		}
		if resp.StatusCode/100 > 2 {
			panic("bad status")
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		item, exist := deploymentsMap[id]
		if !exist {
			c.JSON(http.StatusNotFound, "id not found")
			return
		}

		allocations := []model.MicroserviceAllocation{}
		if err := json.Unmarshal(body, &allocations); err != nil {
			fmt.Println("unmarshal failed")
			panic(err)
		}
		fmt.Println(allocations)
		newDeployments := []*appsv1.Deployment{}
		for _, microservice := range allocations {
			deployment := *item.deployments[microservice.Name]
			for num, container := range microservice.Containers {
				newDeployment := deployment
				newDeployment.Name = fmt.Sprintf("%s-%d", deployment.Name, num)
				newDeployment.Spec.Replicas = utils.NewInt32(1)
				newDeployment.Spec.Template.Spec.NodeSelector = map[string]string{}
				newDeployment.Spec.Template.Spec.NodeSelector["kubernetes.io/hostname"] = container.Loc
				newDeployment.Spec.Selector.MatchLabels["container-allocation-deployment"] = newDeployment.Name
				newDeployment.Spec.Template.Labels["container-allocation-deployment"] = newDeployment.Name
				requests := apiv1.ResourceList{}
				requests[apiv1.ResourceCPU] = *resource.NewMilliQuantity(int64(math.Ceil(container.Cpu)), resource.DecimalSI)
				requests[apiv1.ResourceMemory] = *resource.NewQuantity(utils.Int64(microservice.RequestMemory).MBtoB(), resource.BinarySI)
				newDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = requests
				newDeployment.Spec.Template.Spec.Containers[0].Resources.Limits = nil
				newDeployments = append(newDeployments, &newDeployment)
			}

		}
		for _, deployment := range newDeployments {
			namespace := apiv1.NamespaceDefault
			if deployment.Namespace != "" {
				namespace = deployment.Namespace
			}
			res, err := clientSet.AppsV1().Deployments(namespace).Create(deployment)
			if err != nil {
				fmt.Println(deployment)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			fmt.Println(res)
		}


		c.JSON(200, allocations)
	})
	router.GET("/max-memory", func(c *gin.Context) {
		containerName := c.Query("container")
		query := fmt.Sprintf("container_memory_max_usage_bytes{ container = \"%s\"}", containerName)
		res, err := myProm.Query(query, time.Now())
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("response type is: "+res.Type().String())
		fmt.Println("response string is:"+res.String())
		c.JSON(200, res)
	})

	router.POST("/allocation", DyClientSet(config, &clientSet), func(c *gin.Context) {
		id := c.Query("id")
		item, exist := deploymentsMap[id]
		if !exist {
			c.JSON(http.StatusNotFound, "id not found")
			return
		}
		allocations := []model.MicroserviceAllocation{}
		if err := utils.Bind(c, &allocations); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ret := []*appsv1.Deployment{}
		for _, microservice := range allocations {
			deployment := *item.deployments[microservice.Name]
			for num, container := range microservice.Containers {
				newDeployment := deployment
				newDeployment.Name = fmt.Sprintf("%s-%d", deployment.Name, num)
				newDeployment.Spec.Replicas = utils.NewInt32(1)
				newDeployment.Spec.Template.Spec.NodeSelector = map[string]string{}
				newDeployment.Spec.Template.Spec.NodeSelector["kubernetes.io/hostname"] = container.Loc
				newDeployment.Spec.Selector.MatchLabels["container-allocation-deployment"] = newDeployment.Name
				newDeployment.Spec.Template.Labels["container-allocation-deployment"] = newDeployment.Name
				requests := apiv1.ResourceList{}

				requests[apiv1.ResourceCPU] = *resource.NewMilliQuantity(int64(math.Ceil(container.Cpu)), resource.DecimalSI)
				requests[apiv1.ResourceMemory] = *resource.NewQuantity(utils.Int64(microservice.RequestMemory).MBtoB(), resource.BinarySI)
				newDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = requests
				newDeployment.Spec.Template.Spec.Containers[0].Resources.Limits = nil
				ret = append(ret, &newDeployment)
			}

		}
		for _, deployment := range ret {
			namespace := apiv1.NamespaceDefault
			if deployment.Namespace != "" {
				namespace = deployment.Namespace
			}
			_, err := clientSet.AppsV1().Deployments(namespace).Create(deployment)
			if err != nil {
				fmt.Println(deployment)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
		}


		c.JSON(200, ret)
	})

	router.Run(":8080")
}

func DyClientSet(config *rest.Config, clientSet **kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "token need")
			return
		}
		config.BearerTokenFile = ""
		config.BearerToken = token
		var err error
		if *clientSet, err = kubernetes.NewForConfig(config); err != nil {
			fmt.Println("error occurred in clientSet creating...")
			panic(err.Error())
		}
		//fmt.Printf("clientSet created: %v\n", *clientSet)
	}
}

func DynamicClientSet(config *rest.Config, clientSet **kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "token need")
			return
		}
		authInfo := &api.AuthInfo{
			Token: token,
		}

		//var err error
		//if config, err = rest.InClusterConfig(); err != nil {
		//	fmt.Printf("error occurred in cluster config...")
		//	panic(err.Error())
		//}
		clientConfig := BuildCmdConfig(authInfo, config)
		cfg, err := clientConfig.ClientConfig()
		if err != nil {
			fmt.Printf("error occurred in client config...")
			panic(err.Error())
		}

		*clientSet, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			fmt.Printf("error occurred in client set...")
			panic(err.Error())
		}
		//fmt.Printf("clientSet created: %v\n", *clientSet)

		//config.BearerTokenFile = ""
		//config.BearerToken = token

		//fmt.Println("Before create clientSet...")
		//var err error
		//if clientSet, err = kubernetes.NewForConfig(config); err != nil {
		//	fmt.Println("error occurred in clientSet creating...")
		//	panic(err.Error())
		//}
		//fmt.Printf("Dynamic ClientSet: %v\n", clientSet)
	}
}

func BuildCmdConfig( authInfo *api.AuthInfo, cfg *rest.Config) clientcmd.ClientConfig {
	cmdCfg := api.NewConfig()
	cmdCfg.Clusters[DefaultCmdConfigName] = &api.Cluster{
		Server:                   cfg.Host,
		CertificateAuthority:     cfg.TLSClientConfig.CAFile,
		CertificateAuthorityData: cfg.TLSClientConfig.CAData,
		InsecureSkipTLSVerify:    cfg.TLSClientConfig.Insecure,
	}
	cmdCfg.AuthInfos[DefaultCmdConfigName] = authInfo
	cmdCfg.Contexts[DefaultCmdConfigName] = &api.Context{
		Cluster:  DefaultCmdConfigName,
		AuthInfo: DefaultCmdConfigName,
	}
	cmdCfg.CurrentContext = DefaultCmdConfigName

	return clientcmd.NewDefaultClientConfig(
		*cmdCfg,
		&clientcmd.ConfigOverrides{},
	)
}

func listenAndServe() {
	router := gin.Default()
	router.Use(UseMiddleWares()...)
	handlers.LoadRoutes(router)
	router.Run(":8080")
}

func UseMiddleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc {
		basic.Context(),
		basic.ErrorHandler(),
		middleware.SetInClusterConfig(),
	}
}