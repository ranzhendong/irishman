package datastruck

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/thinkeridea/go-extend/exnet"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//IPPortValidator : Variable declarations
type IPPortValidator struct {
	EMsg string
}

//UpstreamNameValidator : Variable declarations
type UpstreamNameValidator struct {
	EMsg string
}

//PoolNilValidator : Variable declarations
type PoolNilValidator struct {
	EMsg string
}

//Validate : judge ip port values
func (ipv *IPPortValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {

	//regex
	const (
		IP = "^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)$"
	)

	var (
		err error
	)

	defer func() {
		_ = recover()
		if err != nil {
			err = fmt.Errorf("Validate: IPPort Is Valid, EXP: { 192.168.101.59:1080 } ")
		}
	}()

	sep := ":"
	arr := strings.Split(val.String(), sep)

	if !regexp.MustCompile(IP).MatchString(arr[0]) {
		err = fmt.Errorf("Validate: Invalid IP %v ", arr[0])
		return false, err
	}

	if !exnet.HasLocalIPddr(arr[0]) {
		err = fmt.Errorf("Validate: Not Private IP %v ", arr[0])
		return false, err
	}

	a, err := strconv.Atoi(arr[1])
	if int(0) >= a || int(65535) <= a {
		err = fmt.Errorf("Validate: Invalid Port %v ", arr[1])
		return false, err
	}

	return true, nil
}

//Validate : UpstreamName Reserved field filtering
func (unv *UpstreamNameValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {

	//regex
	var (
		err error
		c   Config
	)

	if err = viper.Unmarshal(&c); err != nil {
		err = fmt.Errorf("Validate: Unable To Decode Into Config Struct %v ", err)
		return false, err
	}

	defer func() {
		_ = recover()
		if err != nil {
			err = fmt.Errorf("Validate: Reserved Name is: %v, Your UpstreamName is: %v ", c.Upstream.Reserved, val.String())
		}
	}()

	for _, v := range c.Upstream.Reserved {
		if v == val.String() {
			err = fmt.Errorf("Validate: Reserved Name %v ", v)
			return false, err
		}
	}

	return true, nil
}

//Validate : UpstreamName Reserved field filtering
func (pnv *PoolNilValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {
	var err error

	log.Println(val.Slice(0, 0))
	if val.IsNil() {
		err = fmt.Errorf("Validate: Pool Is None ")
		return false, err
	}

	return true, nil
}
