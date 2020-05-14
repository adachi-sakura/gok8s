package config

import (
	"github.com/buzaiguna/gok8s/utils"
	"io"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
	"strings"
)

const (
	bandwidthFilePath = "/config/result.csv"
)

var (
	PROMETHEUS_HOST	string
	PROMETHEUS_PORT	string
	ALGORITHM_HOST	string
	ALGORITHM_PORT	string
	Bandwidth		int
	InClusterConfig	*rest.Config
)

func init() {
	var err error
	PROMETHEUS_HOST = os.Getenv("PROMETHEUS_SERVICE_HOST")
	PROMETHEUS_PORT = os.Getenv("PROMETHEUS_SERVICE_PORT")
	ALGORITHM_HOST = os.Getenv("CONTAINER_ALLOCATION_SERVICE_HOST")
	ALGORITHM_PORT = os.Getenv("CONTAINER_ALLOCATION_SERVICE_PORT")
	if InClusterConfig, err = rest.InClusterConfig(); err != nil {
		panic(err)
	}
}

func InitBandwidth(line int) {
	Bandwidth = getBandwidth(line)
}

func getBandwidth(line int) int {
	str, err := utils.GetSelectedLineInFile(bandwidthFilePath, line)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	floatBandwidth, err := strconv.ParseFloat(strings.Split(str, ";")[1], 64)
	if err != nil {
		panic(err.Error())
	}
	return int(utils.Int64(floatBandwidth).MBtoKB())
}
