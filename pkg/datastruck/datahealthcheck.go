package datastruck

//HealthCheck
type HealthCheck struct {
	HealthCheckName string   `json:"healthCheckName" validate:"required"`
	Status          string   `json:"status" validate:"required||in=running,stop"`
	CheckProtocol   string   `json:"checkProtocol" validate:"required||in=http,tcp"`
	CheckPath       string   `json:"checkPath" validate:"required"`
	Health          Health   `json:"health" validate:"required"`
	UnHealth        UnHealth `json:"unHealth" validate:"required"`
}

//Health
type Health struct {
	Interval      int   `json:"interval" validate:"required||integer"`
	SuccessTime   int   `json:"successTime" validate:"required||integer"`
	SuccessStatus []int `json:"successStatus" validate:"required||integer||unique"`
}

//UnHealth
type UnHealth struct {
	Interval        int   `json:"interval" validate:"required||integer"`
	FailuresTime    int   `json:"failuresTime" validate:"required||integer"`
	FailuresTimeout int   `json:"failuresTimeout" validate:"required||integer"`
	FailuresStatus  []int `json:"failuresStatus" validate:"required||integer||unique"`
}
