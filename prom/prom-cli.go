package prom

import (
	"fmt"
	"github.com/prometheus/client_golang/api"
	"os"
)

var promCli api.Client

func ConnectProm() api.Client {
	if promCli != nil {
		return promCli
	}
	//todo ns & name parameter
	host, port := os.Getenv("PROMETHEUS_SERVICE_HOST"), os.Getenv("PROMETHEUS_SERVICE_PORT")
	cli, err := getPromCliWithENV(host, port)
	if err != nil {
		panic(err)
	}
	promCli = cli
	return cli
}

func getPromCliWithDNS(ns string, name string, port int) (api.Client, error) {
	cfg := api.Config{
		Address:	fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", name, ns, port),
	}
	c, err := api.NewClient(cfg)
	return c, err
}

func getPromCliWithENV(host string, port string) (api.Client, error) {
	cfg := api.Config{
		Address:	fmt.Sprintf("http://%s:%s", host, port),
	}
	c, err := api.NewClient(cfg)
	return c, err
}