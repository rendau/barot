package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "barot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world!")
	},
}

// Execute - executes root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
