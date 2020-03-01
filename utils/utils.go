package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"strings"
)



type BindError struct {
	reason string
}

func (this BindError) Error() string {
	return this.reason
}

func Bind(c *gin.Context, obj interface{}) error {
	contentType := c.GetHeader("Content-Type")
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	fmt.Println(contentType)
	if contentType == binding.MIMEJSON {
		if err := c.BindJSON(obj); err != nil {
			//panic(err)
			return BindError{
				reason:	fmt.Sprintf("error occurred in json bind due to:\n%s", err.Error()),
			}
		}
	} else if contentType == binding.MIMEYAML {
		yamlBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			//panic(err)
			return BindError{
				reason:	fmt.Sprintf("error occurred in request body get due to:\n%s", err.Error()),
			}
		}
		jsonBody, err := yaml.ToJSON(yamlBody)
		if err != nil {
			//panic(err)
			return BindError{
				reason:	fmt.Sprintf("error occurred in yaml convert due to:\n%s", err.Error()),
			}
		}
		if err := json.Unmarshal(jsonBody, obj); err != nil {
			//panic(err)
			return BindError{
				reason:	fmt.Sprintf("error occurred in yaml-json unmarshal due to:\n%s", err.Error()),
			}
		}

	} else {
		fmt.Println("no content type match")
		return BindError{
			reason: "invalid content type",
		}
	}
	return nil
}
