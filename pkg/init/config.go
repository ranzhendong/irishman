package init

import (
	"github.com/fsnotify/fsnotify"
	MyError "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/spf13/viper"
	"log"
	"os"
)

//Config : load config from config.yaml when irishman start
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
		log.Printf(MyError.ErrorLog(6140), " ", e.Name)
	})

	if pwd, err = os.Getwd(); err != nil {
		os.Exit(1)
		return
	}
	log.Print(MyError.ErrorLog(6141), " ", pwd)

	//Find and read the config and token file
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	return

}
