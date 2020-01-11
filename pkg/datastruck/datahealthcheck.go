package datastruck

//HealthCheck
type HealthCheck struct {
	HealthCheckName string   `json:"healthCheckName" yaml:"healthCheckName" `
	Status          string   `json:"status" yaml:"status" validate:"required||in=running,stop"`
	CheckProtocol   string   `json:"checkProtocol" yaml:"checkProtocol" validate:"required||in=http,tcp"`
	CheckPath       string   `json:"checkPath" yaml:"checkPath" validate:"required"`
	Health          Health   `json:"health" yaml:"health" validate:"required"`
	UnHealth        UnHealth `json:"unhealth" yaml:"unhealth" validate:"required"`
}

//Health
type Health struct {
	Interval      int   `json:"interval" yaml:"interval" validate:"required||integer"`
	SuccessTime   int   `json:"successTime" yaml:"successTime" validate:"required||integer"`
	SuccessStatus []int `json:"successStatus" yaml:"successStatus" validate:"required||integer||unique"`
}

//UnHealth
type UnHealth struct {
	Interval        int   `json:"interval" yaml:"interval" validate:"required||integer"`
	FailuresTime    int   `json:"failuresTime" yaml:"failuresTime" validate:"required||integer"`
	FailuresTimeout int   `json:"failuresTimeout" yaml:"failuresTimeout" validate:"required||integer"`
	FailuresStatus  []int `json:"failuresStatus" yaml:"failuresStatus" validate:"required||integer||unique"`
}
