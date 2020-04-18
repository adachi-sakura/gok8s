package model

import (
	"github.com/buzaiguna/gok8s/utils"
	v1 "k8s.io/api/core/v1"
)

type (
	ResourceQuota struct {
		Cpu_rq_total	int64	`json:"cpu_rq_total"`
		Mem_rq_total	int64	`json:"mem_rq_total"`
	}

	LimitRange struct {
		Cpu_lm	int64	`json:"cpu_lm"`
		Mem_lm	int64	`json:"mem_lm"`
	}

	Node struct {
		Name		string `json:"name"`
		Current_cpu	int64	`json:"current_cpu"`
		Allocatable_cpu int64 `json:"allocatable_cpu"`
		Sum_cpu		int64 `json:"sum_cpu"`
		Current_mem	int64	`json:"current_mem"`
		Allocatable_mem int64 `json:"allocatable_mem"`
		Sum_mem		int64	`json:"sum_mem"`
	}

	Network struct {
		Receive		float64 `json:"receive"`
		Transmit	float64 `json:"transmit"`
	}

	MicroserviceMetrics struct {
		Network
		CpuUsageTime	float64 `json:"cpuUsageTime"`
		CpuTimeTotal	float64	`json:"cpuTimeTotal"`
		HttpRequestsCount	int	`json:"httpRequestCount"`
		MaxMemoryUsage	float64	`json:"maxMemoryUsage"`
	}

	MicroserviceYaml struct {
		Name	string `json:"name"`
		Replicas			int32	`json:"replicas"`
		LeastResponseTime	float64	`json:"leastResponseTime"`
		MicroservicesToInvoke []int	`json:"microservicesToInvoke"`
	}

	MicroservcieData struct {
		MicroserviceMetrics
		MicroserviceYaml
	}

	AlgorithmParameters struct {
		ResourceQuota
		LimitRange
		Nodes 				[]*Node	`json:"nodes"`
		Datas				[]*MicroservcieData	`json:"datas"`
		EntrancePoint		int		`json:"entrancePoint"`
		Bandwidth			int		`json:"bandwidth"`
		TotalTimeRequired	float64	`json:"totalTimeRequired"`
	}

	containerAllocation struct {
		Loc		string `json:"loc"`
		Cpu		float64 `json:"cpu"`
	}

	MicroserviceAllocation struct {
		Name		string `json:"name"`
		Containers	[]containerAllocation `json:"containers"`
		RequestMemory float64 `json:"requestMemory"`
	}


)

func NewMicroserviceYaml() *MicroserviceYaml {
	return &MicroserviceYaml{
		MicroservicesToInvoke: []int{},
	}
}

func (this *MicroserviceYaml) SetLeastResponseTime(str string) error {
	if str == "" {
		this.LeastResponseTime = float64(utils.INT_MAX)
		return nil
	}
	var err error
	this.LeastResponseTime, err = utils.Float64(str)
	return err
}

func NewLimitRange(limitRanges *v1.LimitRangeList) LimitRange {
	lm := LimitRange{utils.INT64_MAX, utils.INT64_MAX}
	for _, limitRange := range limitRanges.Items {
		for _, item := range limitRange.Spec.Limits {
			if item.Type != v1.LimitTypeContainer {
				continue
			}
			if item.Max == nil {
				continue
			}
			if maxCpu, exists := item.Max[v1.ResourceCPU]; exists {
				lm.Cpu_lm = utils.Int64Min(lm.Cpu_lm, maxCpu.MilliValue())
			}
			if maxMem, exists := item.Max[v1.ResourceMemory]; exists {
				lm.Mem_lm = utils.Int64Min(lm.Mem_lm, maxMem.Value())
			}
		}
	}
	return lm
}

func NewResourceQuota(resourceQuotas *v1.ResourceQuotaList) ResourceQuota {
	rq := ResourceQuota{utils.INT64_MAX, utils.INT64_MAX}
	for _, resourceQuota := range resourceQuotas.Items {
		if resourceQuota.Spec.Hard == nil {
			continue
		}
		if maxCpu, exists := resourceQuota.Spec.Hard[v1.ResourceCPU]; exists {
			rq.Cpu_rq_total = utils.Int64Min(rq.Cpu_rq_total, maxCpu.MilliValue())
		}
		if maxMem, exists := resourceQuota.Spec.Hard[v1.ResourceMemory]; exists {
			rq.Mem_rq_total = utils.Int64Min(rq.Mem_rq_total, maxMem.Value())
		}
	}
	return rq
}

func (this *AlgorithmParameters) SetTotalTimeRequired(str string) error {
	if str == "" {
		this.TotalTimeRequired = float64(utils.INT_MAX)
		return nil
	}
	var err error
	this.TotalTimeRequired, err = utils.Float64(str)
	return err
}