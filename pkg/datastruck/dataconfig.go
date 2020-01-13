package datastruck

import (
	"github.com/spf13/viper"
)

//config.yaml
type Config struct {
	Server      server      `yaml:"server"`
	Etcd        etcd        `yaml:"etcd"`
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
		return
	}
	return nil
}