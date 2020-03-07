package main
import (
	"context"
	"fmt"
	prom_cli "github.com/buzaiguna/gok8s/prom-cli"
	"github.com/buzaiguna/gok8s/utils"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"net/http"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
)

const (
	DefaultCmdConfigName = "kubernetes"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("initial config token: %s\ntokenfile: %s", config.BearerToken, config.BearerTokenFile)
	clientSet := &kubernetes.Clientset{}


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
		cli := prom_cli.ConnectProm()
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
