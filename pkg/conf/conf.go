package conf

import (
	"fmt"
	"log"

	"github.com/creasty/defaults"
	"github.com/cxfksword/go-docker-skeleton/pkg/utils"
	"github.com/spf13/viper"
)

var currentAppName string
var currentPath string

func Init(appName string, path string) {
	currentAppName = appName
	currentPath = path

	if path == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(defaultConfPath(appName))
	} else {
		if !utils.Exists(path) {
			log.Panicf("Config file not exists. file: %s \n", path)
		}

		viper.SetConfigFile(path)
	}

	err := defaults.Set(App)
	if err != nil { // Handle errors reading the config file
		log.Fatalf("Config file read failed. error: %s \n", err)
	}

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Printf("Config file read failed. error: %s \n", err)
		return
	}

	err = viper.Unmarshal(App)
	if err != nil {
		fmt.Printf("Config unable to decode into struct, %v\n", err)
	}

}

func Reload() {
	Init(currentAppName, currentPath)
}
