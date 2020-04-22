package appclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buzaiguna/gok8s/apperror"
	cfg "github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/model"
	"github.com/buzaiguna/gok8s/utils"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
)

const (
	algorithmUrlBase = "http://%s:%s/algorithm"
)

var algorithmUrl string

type AlgorithmClient struct {
	client	*http.Client
}

func NewAlgorithmClient() *AlgorithmClient {
	return &AlgorithmClient{
		client:	&http.Client{},
	}
}

func init() {
	algorithmUrl = fmt.Sprintf(algorithmUrlBase, cfg.ALGORITHM_HOST, cfg.ALGORITHM_PORT)
}

func (cli *AlgorithmClient) GetAllocations(params *model.AlgorithmParameters) ([]model.MicroserviceAllocation, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := cli.DoRequest(http.MethodGet, algorithmUrl, jsonBytes)
	defer resp.Body.Close()
	if utils.BadResponse(resp) {
		return nil, buildBadResponseError(resp)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	allocations := []model.MicroserviceAllocation{}
	if err := json.Unmarshal(body, &allocations); err != nil {
		return nil, err
	}
	return allocations, nil
}

func (cli *AlgorithmClient) DoRequest(method string, url string, body []byte) *http.Response {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	if body != nil {
		request.Header.Set("Content-Type", binding.MIMEJSON)
	}

	resp, err := cli.Do(request)
	if err != nil {
		panic(err)
	}
	return resp
}

func (cli *AlgorithmClient) Do(request *http.Request) (*http.Response, error) {
	return cli.client.Do(request)
}

func buildBadResponseError(resp *http.Response) error {
	body, _ := ioutil.ReadAll(resp.Body)
	return apperror.NewInternalServerError(string(body[:]))
}