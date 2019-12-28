package datastruck

import "encoding/json"

//config.yaml
type Config struct {
	Etcd etcd `yaml:"etcd"`
}

type etcd struct {
	Host   string `yaml:"host"`
	Format string `yaml:"format"`
}

//Upstream
type Upstream struct {
	UpstreamName string   `json:"upstreamName"`
	Algorithms   string   `json:"algorithms"`
	Pool         []Server `json:"pool"`
}

type Server struct {
	IpPort string      `json:"ipPort"`
	Status string      `json:"status"`
	Weight json.Number `json:"weight"`
}
