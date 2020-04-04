package config

import (
	"github.com/buzaiguna/gok8s/utils"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	bandwidthFilePath = "/config/result.csv"
	lineRemoteVMTCP = 5
)

var (
	PROMETHEUS_HOST	string
	PROMETHEUS_PORT	string
	ALGORITHM_HOST	string
	ALGORITHM_PORT	string
	Bandwidth		int
)

func init() {
	PROMETHEUS_HOST = os.Getenv("PROMETHEUS_SERVICE_HOST")
	PROMETHEUS_PORT = os.Getenv("PROMETHEUS_SERVICE_PORT")
	ALGORITHM_HOST = os.Getenv("CONTAINER_ALLOCATION_SERVICE_HOST")
	ALGORITHM_PORT = os.Getenv("CONTAINER_ALLOCATION_SERVICE_PORT")
	Bandwidth = getBandwidth()
}

func getBandwidth() int {
	str, err := utils.GetSelectedLineInFile(bandwidthFilePath, lineRemoteVMTCP)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	floatBandwidth, err := strconv.ParseFloat(strings.Split(str, ";")[1], 64)
	if err != nil {
		panic(err.Error())
	}
	return int(utils.Int64(floatBandwidth).MBtoKB())
}
