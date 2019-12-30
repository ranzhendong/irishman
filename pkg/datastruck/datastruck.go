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
	UpstreamName string   `json:"upstreamName" validate:"required"`
	Algorithms   string   `json:"algorithms" validate:"required||in=ip-hex,round-robin"`
	Pool         []Server `json:"pool" validate:"required"`
}

type Server struct {
	IpPort string      `json:"ipPort" validate:"required"`
	Status string      `json:"status" validate:"required"`
	Weight json.Number `json:"weight" validate:"required"`
}
