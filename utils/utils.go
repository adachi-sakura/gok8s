package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
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

func DecodeK8SResources(c *gin.Context) ([]runtime.Object, error) {
	yamlFiles, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//panic(err)
		return nil, BindError{
			reason:	fmt.Sprintf("error occurred in request body get due to:\n%s", err.Error()),
		}
	}
	return parseK8SYaml(yamlFiles), nil
}

func parseK8SYaml(fileR []byte) []runtime.Object {
	filesAsString := string(fileR[:])
	sepYamlFiles := strings.Split(filesAsString, "---")
	retObj := []runtime.Object{}
	for _, file := range sepYamlFiles {
		if file == "\n" || file == "" {
			continue
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode([]byte(file), nil, nil)
		fmt.Println(groupVersionKind)
		if err != nil {
			fmt.Printf("error occurred when decoding yaml file\n %s ", err.Error())
			continue
		}
		retObj = append(retObj, obj)
	}
	return retObj
}

func TrimSpace(strs []string) []string {
	ret := []string{}
	for _, str := range strs {
		ret = append(ret, strings.TrimSpace(str))
	}
	return ret
}
