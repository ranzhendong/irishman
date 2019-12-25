package initconfig

import (
	"datastruck"
	"github.com/spf13/viper"
	"log"
)

func ConfigAnalysis() (err error, c datastruck.Config) {
	// init viper config
	if err = Config(); err != nil {
		log.Println(err)
		return
	}

	// Unmarshal to struck
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[Main] Unable To Decode Into Config Struct, %v", err)
		return
	}
	return
}
