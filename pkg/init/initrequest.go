package init

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

//InitializeBody : config initialize
func InitializeBody(rBody io.Reader) (jsonObj interface{}, err error) {
	var (
		body []byte
	)

	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		err = fmt.Errorf("Read Body ERR: %v ", err)
		return
	}

	// if the body can be turn to interface
	if err = json.Unmarshal(body, &jsonObj); err != nil {
		err = fmt.Errorf("Unmarshal Body ERR: %v", err)
		return
	}

	return
}
