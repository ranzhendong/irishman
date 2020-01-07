package reconstruct

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Reserved []string `yaml:"reserved"`
}

type IpPortValidator struct {
	EMsg string
}

type UpstreamNameValidator struct {
	EMsg string
}

//ip port vaildator
func (self *IpPortValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {

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
			//log.Printf("[Validate]: IP, Port or : Is Valid, exp: { 192.168.101.59:1080 }")
			err = fmt.Errorf("[Validate]: IP, Port or : Is Valid, exp: { 192.168.101.59:1080 }")
		}
	}()

	sep := ":"
	arr := strings.Split(val.String(), sep)

	if !regexp.MustCompile(IP).MatchString(arr[0]) {
		err = fmt.Errorf("[Validate]: IP Invalid , IP: %v", arr[0])
		//log.Printf("[Validate]: IP Invalid , IP: %v", arr[0])
		return false, err
	}

	a, err := strconv.Atoi(arr[1])
	if int(1024) >= a || int(65535) <= a {
		err = fmt.Errorf("[Validate]: Port Invalid , Port: %v", arr[1])
		//log.Printf("[Validate]: Port Invalid , Port: %v", arr[1])
		return false, err
	}

	return true, nil
}

//UpstreamName Reserved field filtering
func (self *UpstreamNameValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {

	//regex
	var (
		err error
		c   Config
	)

	if err = viper.Unmarshal(&c); err != nil {
		err = fmt.Errorf("[Validate] Unable To Decode Into Config Struct, %v", err)
		return false, err
	}

	defer func() {
		_ = recover()
		if err != nil {
			err = fmt.Errorf("[Validate]: Reserved Name is: %v, Your UpstreamName is: %v", c.Reserved, val.String())
		}
	}()

	for _, v := range c.Reserved {
		if v == val.String() {
			err = fmt.Errorf("[Validate] Reserved Name, %v", v)
			return false, err
		}
	}

	return true, err
}
