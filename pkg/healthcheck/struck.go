package healthcheck

//GoroutinesHealthCheck : write HC template to nutsDB
type GoroutinesHealthCheck struct {
	CheckProtocol string `json:"checkProtocol"`
	CheckPath     string `json:"checkPath"`
	Health        struct {
		Interval       int   `json:"interval"`
		SuccessTime    int   `json:"successTime"`
		SuccessTimeout int   `json:"successTimeout"`
		SuccessStatus  []int `json:"successStatus"`
	} `json:"health"`
	UnHealth struct {
		Interval        int   `json:"interval"`
		FailuresTime    int   `json:"failuresTime"`
		FailuresTimeout int   `json:"failuresTimeout"`
		FailuresStatus  []int `json:"failuresStatus"`
	} `json:"unhealth"`
}

//Upstream: write upstream server up&down to nutsDB
type Upstream struct {
	UpstreamName string `json:"upstreamName"`
	Pool         []struct {
		IPPort string `json:"ipPort"`
		Status string `json:"status"`
		Weight int    `json:"weight"`
	} `json:"pool"`
}

//Config: for methods
type TConfig struct {
	UpstreamEtcPrefix    string
	HealthCheckEtcPrefix string
	TagUp                string
	TagDown              string
	TagSuccessCode       string
	TagFailureCode       string
}
