package datastruck

import (
	ErrH "errorhandle"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

//config.yaml
type Config struct {
	Server server `yaml:"server"`
	Etcd   etcd   `yaml:"etcd"`
	//Resource    resource    `yaml:"resource"`
	NutsDB      nutsDB      `yaml:"nutsdb"`
	Upstream    upstream    `yaml:"upstream"`
	HealthCheck healthcheck `yaml:"healthcheck"`
}

type server struct {
	Bind         string `yaml:"bind"`
	ReadTimeout  int    `yaml:"readtimeout"`
	WriteTimeout int    `yaml:"writetimeout"`
}

type etcd struct {
	Host    string `yaml:"host"`
	Format  string `yaml:"format"`
	Timeout int    `yaml:"timeout"`
}

type resource struct {
	UpstreamList string `yaml:"upstreamlist"`
	UpList       string `yaml:"uplist"`
	DownList     string `yaml:"downlist"`
}

type nutsDB struct {
	Path string `yaml:"path"`
	Tag  struct {
		Up           string `yaml:"up"`
		Down         string `yaml:"down"`
		UpstreamList string `yaml:"upstreamlist"`
	} `yaml:"tag"`
}

type upstream struct {
	EtcdPrefix string   `yaml:"etcdprefix"`
	Reserved   []string `yaml:"reserved"`
}

type healthcheck struct {
	EtcdPrefix string      `yaml:"etcdprefix"`
	Template   HealthCheck `yaml:"template"`
}

func (self *Config) Config() (err error) {
	if err = viper.Unmarshal(&self); err != nil {
		log.Printf(ErrH.ErrorLog(0142), fmt.Sprintf("%v", err))
		return
	}
	return nil
}
