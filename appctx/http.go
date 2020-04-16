package appctx

import (
	"context"
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
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