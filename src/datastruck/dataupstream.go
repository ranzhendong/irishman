package datastruck

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/smokezl/govalidators"
	"log"
)

//Upstream, for put post
type Upstream struct {
	UpstreamName string   `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string   `json:"algorithms" validate:"required||in=ip-hex,round-robin"`
	Pool         []Server `json:"pool" validate:"required"`
}

//upstream server, for put post
type Server struct {
	IpPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"required||in=up,down,nohc"`
	Weight int    `json:"weight" validate:"required||integer"`
}

//Upstream, for get
type GetUpstream struct {
	UpstreamName string `json:"upstreamName" validate:"required"`
}

//Upstream, for patch
type PatchUpstream struct {
	UpstreamName string        `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string        `json:"algorithms" validate:"in=ip-hex,round-robin"`
	Pool         []PatchServer `json:"pool"`
}

//upstream server, for patch
type PatchServer struct {
	IpPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"in=up,down,nohc"`
	Weight int    `json:"weight" validate:"integer"`
}

//Upstream, for delete
type DeleteUpstream struct {
	UpstreamName string         `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string         `json:"algorithms"`
	Pool         []DeleteServer `json:"pool"`
}

//upstream server, for delete
type DeleteServer struct {
	IpPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status"`
	Weight int    `json:"weight"`
}

func (self *Upstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &self); err != nil {
		return
	}

	//new filter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &IpPortValidator{},
		"upstreamName": &UpstreamNameValidator{},
	})

	//if not match
	if err := validator.Validate(self); err != nil {
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

func (self *GetUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &self); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(self); err != nil {
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

func (self *PatchUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &self); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &IpPortValidator{},
		"upstreamName": &UpstreamNameValidator{},
		"poolNil":      &PoolNilValidator{},
	})

	//if not match
	if err := validator.Validate(self); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

func (self *DeleteUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &self); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &IpPortValidator{},
		"upstreamName": &UpstreamNameValidator{},
	})

	//if not match
	if err := validator.Validate(self); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}
