package datastruck

import (
	"fmt"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/spf13/viper"
	"log"
)

//Config is struck of config.yaml
type Config struct {
	Server      server      `yaml:"server"`
	PubSub      pubSub      `yaml:"pubsub"`
	Metrics     metrics     `yaml:"metrics"`
	Etcd        etcd        `yaml:"etcd"`
	Resource    resource    `yaml:"resource"`
	NutsDB      nutsDB      `yaml:"nutsdb"`
	Upstream    upstream    `yaml:"upstream"`
	HealthCheck healthcheck `yaml:"healthcheck"`
}

type server struct {
	Bind         string `yaml:"bind"`
	ReadTimeout  int    `yaml:"readtimeout"`
	WriteTimeout int    `yaml:"writetimeout"`
}

type pubSub struct {
	Bind         string `yaml:"bind"`
	ReadTimeout  int    `yaml:"readtimeout"`
	WriteTimeout int    `yaml:"writetimeout"`
}

type metrics struct {
	Bind         string `yaml:"bind"`
	Interval     int    `yaml:"interval"`
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
		SuccessCode  string `yaml:"successcode"`
		FailureCode  string `yaml:"failurecode"`
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

//Config: for methods
type TConfig struct {
	UpstreamEtcPrefix    string
	HealthCheckEtcPrefix string
	TagUp                string
	TagDown              string
	TagSuccessCode       string
	TagFailureCode       string
}

//Config : Unmarshal the config
func (c *Config) Config() (err error) {
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf(MyERR.ErrorLog(0142), fmt.Sprintf("%v", err))
		return
	}
	return nil
}

//TC: for functions: SeparateUpstreamToNuts, HealthCheckTemplateToNuts
func (c *Config) TC() TConfig {

	//specify the explicit value, just for values not exist.
	return TConfig{
		UpstreamEtcPrefix:    c.Upstream.EtcdPrefix,
		HealthCheckEtcPrefix: c.HealthCheck.EtcdPrefix,
		TagUp:                c.NutsDB.Tag.Up,
		TagDown:              c.NutsDB.Tag.Down,
		TagSuccessCode:       c.NutsDB.Tag.SuccessCode,
		TagFailureCode:       c.NutsDB.Tag.FailureCode,
	}
}
