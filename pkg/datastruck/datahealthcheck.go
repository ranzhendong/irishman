package datastruck

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/smokezl/govalidators"
	"log"
)

//HealthCheck for template and put
type HealthCheck struct {
	HealthCheckName string   `json:"healthCheckName" yaml:"healthCheckName" validate:"required"`
	Status          string   `json:"status" yaml:"status" validate:"required||in=running,stop"`
	CheckProtocol   string   `json:"checkProtocol" yaml:"checkProtocol" validate:"required||in=http,tcp"`
	CheckPath       string   `json:"checkPath" yaml:"checkPath" validate:"required"`
	Health          Health   `json:"health" yaml:"health" validate:"required"`
	UnHealth        UnHealth `json:"unhealth" yaml:"unhealth" validate:"required"`
}

//Health : template and put Health
type Health struct {
	Interval       int   `json:"interval" yaml:"interval" validate:"required||integer"`
	SuccessTime    int   `json:"successTime" yaml:"successTime" validate:"required||integer"`
	SuccessTimeout int   `json:"successTimeout" yaml:"successTimeout" validate:"required||integer"`
	SuccessStatus  []int `json:"successStatus" yaml:"successStatus" validate:"required||unique||array"`
}

//UnHealth : template and put UnHealth
type UnHealth struct {
	Interval        int   `json:"interval" yaml:"interval" validate:"required||integer"`
	FailuresTime    int   `json:"failuresTime" yaml:"failuresTime" validate:"required||integer"`
	FailuresTimeout int   `json:"failuresTimeout" yaml:"failuresTimeout" validate:"required||integer"`
	FailuresStatus  []int `json:"failuresStatus" yaml:"failuresStatus" validate:"required||unique||array"`
}

//GetHealthCheck : get method for HealthCheck
type GetHealthCheck struct {
	HealthCheckName string `json:"healthCheckName" validate:"required"`
}

//PatchHealthCheck : patch method for HealthCheck
type PatchHealthCheck struct {
	HealthCheckName string        `json:"healthCheckName" validate:"required"`
	Status          string        `json:"status" validate:"required||in=running,stop"`
	CheckProtocol   string        `json:"checkProtocol"  validate:"required||in=http,tcp"`
	CheckPath       string        `json:"checkPath" validate:"required"`
	Health          PatchHealth   `json:"health" `
	UnHealth        PatchUnHealth `json:"unhealth"`
}

//PatchHealth : patch method for HealthCheck's Health
type PatchHealth struct {
	Interval       int   `json:"interval" validate:"integer"`
	SuccessTime    int   `json:"successTime" validate:"integer"`
	SuccessTimeout int   `json:"successTimeout" validate:"integer"`
	SuccessStatus  []int `json:"successStatus" validate:"unique||array"`
}

//PatchUnHealth : patch method for HealthCheck's UnHealth
type PatchUnHealth struct {
	Interval        int   `json:"interval" validate:"integer"`
	FailuresTime    int   `json:"failuresTime" validate:"integer"`
	FailuresTimeout int   `json:"failuresTimeout" validate:"integer"`
	FailuresStatus  []int `json:"failuresStatus" validate:"unique||array"`
}

//DeleteHealthCheck : patch delete for HealthCheck
type DeleteHealthCheck struct {
	HealthCheckName string `json:"healthCheckName" validate:"required"`
	Health          struct {
		SuccessStatus []int `json:"successStatus" validate:"unique||array"`
	} `json:"health"`
	UnHealth struct {
		FailuresStatus []int `json:"failuresStatus" validate:"unique||array"`
	} `json:"unhealth"`
}

//JudgeValidator : judge the HealthCheck template values
func (hc *HealthCheck) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &hc); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(hc); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the GetHealthCheck template values
func (ghc *GetHealthCheck) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &ghc); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(ghc); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the PatchHealthCheck template values
func (phc *PatchHealthCheck) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &phc); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(phc); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}

//JudgeValidator : judge the DeleteHealthCheck template values
func (dhc *DeleteHealthCheck) JudgeValidator(jsonObj interface{}) (err error) {
	//turn map to struck
	if err = mapstructure.Decode(jsonObj, &dhc); err != nil {
		return
	}

	//judge parameter
	validator := govalidators.New()
	validator.SetValidators(map[string]interface{}{})

	//if not match
	if err := validator.Validate(dhc); err != nil {
		log.Println(err)
		err := fmt.Errorf("%v", err[0])
		return err
	}

	return nil
}
