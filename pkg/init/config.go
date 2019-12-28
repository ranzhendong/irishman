package init

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func Config() (err error) {
	var (
		pwd string
	)

	viper.New()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")

	//watch the config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("[LoadConfig] Config file changed:", e.Name)
	})

	if pwd, err = os.Getwd(); err != nil {
		os.Exit(1)
		return
	}
	log.Println("[LoadConfig] lrishMan Is Running, Execute Path", pwd)

	//Find and read the config and token file
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("[LoadConfig] Fatal Error Config File: %s \n", err)
		err = fmt.Errorf("[LoadConfig] Fatal Error Config File: %s \n", err)
		return
	}
	return

}
