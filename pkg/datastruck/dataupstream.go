package datastruck

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/smokezl/govalidators"
	"log"
)

//Upstream : for put post
type Upstream struct {
	UpstreamName string   `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string   `json:"algorithms" validate:"required||in=ip-hex,round-robin"`
	Pool         []Server `json:"pool" validate:"required"`
}

//Server : upstream server, for put post
type Server struct {
	IPPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"required||in=up,down,nohc"`
	Weight int    `json:"weight" validate:"required||integer"`
}

//GetUpstream : Upstream, for get
type GetUpstream struct {
	UpstreamName string `json:"upstreamName" validate:"required"`
}

//PatchUpstream : Upstream, for patch
type PatchUpstream struct {
	UpstreamName string        `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string        `json:"algorithms" validate:"in=ip-hex,round-robin"`
	Pool         []PatchServer `json:"pool"`
}

//PatchServer : upstream server, for patch
type PatchServer struct {
	IPPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status" validate:"in=up,down,nohc"`
	Weight int    `json:"weight" validate:"integer"`
}

//DeleteUpstream : Upstream, for delete
type DeleteUpstream struct {
	UpstreamName string         `json:"upstreamName" validate:"required||upstreamName"`
	Algorithms   string         `json:"algorithms"`
	Pool         []DeleteServer `json:"pool"`
}

//DeleteServer : upstream server, for delete
type DeleteServer struct {
	IPPort string `json:"ipPort" validate:"required||unique||ipPort"`
	Status string `json:"status"`
	Weight int    `json:"weight"`
}

//JudgeValidator : judge the Upstream template values
func (us *Upstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &us); err != nil {
		return
	}

	//new filter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort": &IPPortValidator{},
	})

	//if not match
	if err := validator.Validate(us); err != nil {
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the GetUpstream template values
func (gus *GetUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &gus); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(gus); err != nil {
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the PatchUpstream template values
func (pus *PatchUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &pus); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &IPPortValidator{},
		"upstreamName": &UpstreamNameValidator{},
		"poolNil":      &PoolNilValidator{},
	})

	//if not match
	if err := validator.Validate(pus); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the DeleteUpstream template values
func (dus *DeleteUpstream) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &dus); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{
		"ipPort":       &IPPortValidator{},
		"upstreamName": &UpstreamNameValidator{},
	})

	//if not match
	if err := validator.Validate(dus); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}
