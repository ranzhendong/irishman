package init

import (
	"datastruck"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func InitializeBody(rBody io.Reader) (err error, u datastruck.Upstream) {
	var (
		body []byte
	)
	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		log.Printf("[InitCheck] Read Body ERR: %v\n", err)
		err = fmt.Errorf("[InitCheck] Read Body ERR: %v\n", err)
		return
	}

	// if the body can be turn to json
	if err = json.Unmarshal(body, &u); err != nil {
		log.Printf("[InitCheck] Unmarshal Body ERR: %v", err)
		err = fmt.Errorf("[InitCheck] Unmarshal Body ERR: %v", err)
		return
	}

	//// log the parameter
	//if parameter, err := json.Marshal(u); err == nil {
	//	log.Printf("[InitCheck] The Request Body: %v", string(parameter))
	//}
	return
}
