package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

// Execute - executes root command
func Execute() {
	loadConf()
	fmt.Println("Hello world!")
}

func loadConf() {
	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath != "" {
		viper.SetConfigFile(confFilePath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	// env vars are in priority
	viper.AutomaticEnv()
}
