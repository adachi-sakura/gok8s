package model

type (
	Cost struct {
		Name		string `json:"name"`
		BasePrice 	string `json:"basePrice"`
		UnitPrice	string `json:"unitPrice"`
	}

	GamspNodeInfo struct {
		Name      string	`json:"name"`
		MilliCore int64		`json:"milliCore"`
		Mem       int64		`json:"mem"`
		BasePrice float64 	`json:"basePrice"`
		UnitPrice float64 	`json:"unitPrice"`
	}

	GamspParameters struct {
		Nodes	[]*GamspNodeInfo	`json:"nodes"`
	}

	pod struct {
		Loc		string `json:"loc"`
		Cpu		float64	`json:"cpu"`
	}

	GamspAllocation struct {
		Pods	[]pod	`json:"pods"`
	}
)
