package datastruck

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
	IpPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"required||in=up,down,nohc"`
	Weight int    `json:"weight" validate:"required||integer"`
}
