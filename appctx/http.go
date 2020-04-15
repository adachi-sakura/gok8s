package appctx

import (
	"context"
	"encoding/json"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/buzaiguna/gok8s/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Param(ctx context.Context, key string) string {
	c := GinContext(ctx)
	return c.Param(key)
}

func Query(ctx context.Context, key string) string {
	c := GinContext(ctx)
	return c.Query(key)
}

func Bind(ctx context.Context, obj interface{}) error {
	c := GinContext(ctx)
	var err error
	switch c.ContentType() {
	case binding.MIMEJSON:
		err =  c.BindJSON(obj)
	case binding.MIMEYAML:
		err = c.BindYAML(obj)
	default:
		err = c.Bind(obj)
	}
	if err != nil {
		return apperror.NewInvalidRequestBodyError(err)
	}
	return nil
}

func BindK8SYaml(ctx context.Context, obj interface{}) error {
	c := GinContext(ctx)
	jsonBody := getYamlBody(c).yamlToJson()
	if err := json.Unmarshal(jsonBody, obj); err != nil {
		return apperror.NewInvalidRequestBodyError(err)
	}
	return nil
}

func DecodeMultiK8SResource(ctx context.Context) []runtime.Object {
	c := GinContext(ctx)
	yamlFiles := getYamlBody(c)
	return utils.ParseK8SYaml(yamlFiles)
}

type requestBody []byte

func getYamlBody(c *gin.Context) requestBody {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func (body requestBody) yamlToJson() []byte {
	jsonBody, err := yaml.ToJSON(body)
	if err != nil {
		panic(err)
	}
	return jsonBody
}