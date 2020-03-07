package prom_cli

import (
	"fmt"
	"github.com/prometheus/client_golang/api"
)

var promCli api.Client

func ConnectProm() api.Client {
	if promCli != nil {
		return promCli
	}
	//todo ns & name parameter
	ns := "monitoring"
	name := "prometheus"
	port := 9090
	cli, err := getPromCli(ns, name, port)
	if err != nil {
		panic(err)
	}
	promCli = cli
	return cli
}

func getPromCli(ns string, name string, port int) (api.Client, error) {
	cfg := api.Config{
		Address:	fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", name, ns, port),
	}
	c, err := api.NewClient(cfg)
	return c, err
}