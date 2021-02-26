package Config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitViper() {
	viper.AddConfigPath("resource/.")
	viper.SetConfigName("application")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("cfg file error: %s\n", err)
		os.Exit(1)
	}
}
