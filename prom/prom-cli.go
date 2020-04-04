package prom

import (
	"context"
	"fmt"
	"github.com/buzaiguna/gok8s/config"
	"github.com/prometheus/client_golang/api"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

var promCli api.Client

func ConnectProm() api.Client {
	if promCli != nil {
		return promCli
	}
	//todo ns & name parameter
	host, port := config.PROMETHEUS_HOST, config.PROMETHEUS_PORT
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

func Query(query string, t time.Time) (model.Value, error) {
	cli := ConnectProm()
	res, _, err := prom.NewAPI(cli).Query(context.Background(), query, t)
	return res, err
}