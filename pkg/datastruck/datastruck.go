package datastruck

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/smokezl/govalidators"
	"log"
	"reconstruct"
)

//config.yaml
type Config struct {
	Etcd     etcd     `yaml:"etcd"`
	Reserved []string `yaml:"reserved"`
}

type etcd struct {
	Host    string `yaml:"host"`
	Format  string `yaml:"format"`
	Timeout int    `yaml:"timeout"`
}

//Upstream, for put post
type Upstream struct {
	UpstreamName string   `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string   `json:"algorithms" validate:"required||in=ip-hex,round-robin"`
	Pool         []Server `json:"pool" validate:"required"`
}

//Upstream, for get
type GetUpstream struct {
	UpstreamName string `json:"upstreamName" validate:"required"`
}

//upstream server
type Server struct {
	IpPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"required||in=up,down,nohc"`
	Weight int    `json:"weight" validate:"required||integer"`
}

func (u *Upstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err := mapstructure.Decode(jsonObj, &u); err != nil {
		fmt.Println(err)
	}

	//new filter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &reconstruct.IpPortValidator{},
		"upstreamName": &reconstruct.UpstreamNameValidator{},
	})

	//if not match
	if err := validator.Validate(u); err != nil {
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

func (gu *GetUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err := mapstructure.Decode(jsonObj, &gu); err != nil {
		fmt.Println(err)
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(gu); err != nil {
		log.Println(err)
		err := fmt.Errorf("ERR: %v", err)
		return err
	}

	return nil
}
