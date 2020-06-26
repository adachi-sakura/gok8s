package GAMSP

import (
	"encoding/json"
	"github.com/buzaiguna/gok8s/model"
	"github.com/buzaiguna/gok8s/utils"
	"testing"
)

func TestBuildMicroserviceDeployments(t *testing.T) {
	allocationsBytes := `[
{
        "pods": [
                {
                        "cpu": 184,
                        "loc": "tongji-k8s-worker1" 
                } 
        ] 
}, 

{
        "pods": [
                {
                        "cpu": 210,
                        "loc": "tongji-k8s-worker1" 
                },
                {
                        "cpu": 939,
                        "loc": "tongji-k8s-worker3" 
                } 
        ] 
}, 

{
        "pods": [
                {
                        "cpu": 197,
                        "loc": "tongji-k8s-worker3" 
                } 
        ] 
}, 

{
        "pods": [
                {
                        "cpu": 663,
                        "loc": "tongji-k8s-worker1" 
                },
                {
                        "cpu": 136,
                        "loc": "tongji-k8s-worker3" 
                } 
        ] 
}, 

{
        "pods": [
                {
                        "cpu": 475,
                        "loc": "tongji-k8s-worker1" 
                },
                {
                        "cpu": 194,
                        "loc": "tongji-k8s-worker3" 
                } 
        ] 
}
]`
	allocations := []model.GamspAllocation{}
	if err := json.Unmarshal([]byte(allocationsBytes), &allocations); err != nil {
		t.Error(err)
	}
	deployments := buildMicroserviceDeployments(allocations)
	t.Log(utils.Stringify(deployments))

}