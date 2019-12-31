package init

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func InitializeBody(rBody io.Reader) (err error, jsonObj interface{}) {
	var (
		body []byte
	)

	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		log.Printf("[InitCheck] Read Body ERR: %v\n", err)
		err = fmt.Errorf("[InitCheck] Read Body ERR: %v\n", err)
		return
	}

	// if the body can be turn to interface
	if err = json.Unmarshal(body, &jsonObj); err != nil {
		log.Printf("[InitCheck] Unmarshal Body ERR: %v", err)
		err = fmt.Errorf("[InitCheck] Unmarshal Body ERR: %v", err)
		return
	}

	return
}
