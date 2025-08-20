package Config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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
