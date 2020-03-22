package model

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
		Current_cpu	int64	`json:"current_cpu"`
		Allocatable_cpu int64 `json:"allocatable_cpu"`
		Sum_cpu		int64 `json:"sum_cpu"`
		Current_mem	int64	`json:"current_mem"`
		Allocatable_mem int64 `json:"allocatable_mem"`
		Sum_mem		int64	`json:"sum_mem"`
	}

	Network struct {
		receive		float64 `json:"receive"`
		transmit	float64 `json:"transmit"`
	}

	MicroserviceMetrics struct {
		Network
		CpuUsageTime	float64 `json:"cpuUsageTime"`
		CpuTimeTotal	float64	`json:"cpuTimeTotal"`
		HttpRequestsCount	int	`json:"httpRequestCount"`
		MaxMemoryUsage	float64	`json:"maxMemoryUsage"`
	}

	MicroserviceYaml struct {
		Replicas			int	`json:"replicas"`
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
)