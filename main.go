package main
import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"net/http"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "PONG")
	})
	router.GET("/pods", func(c *gin.Context) {
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, pods)

	})
	router.GET("/first-pod", func(c *gin.Context) {
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, pods.Items[0])

	})
	router.POST("/deployments", func(c * gin.Context) {
		deployment := &appsv1.Deployment{}
		if err := c.Bind(deployment); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		result, err := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Create(deployment)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusCreated, result)
	})
	router.Run(":8080")
}
