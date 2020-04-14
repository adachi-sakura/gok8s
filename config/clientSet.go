package config

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)


func NewClient(config *rest.Config) *kubernetes.Clientset {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientSet
}

