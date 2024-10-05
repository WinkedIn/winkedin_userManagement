package services

import (
	"flag"
	"github.com/spf13/viper"
	"sync"
)

var (
	ConfigFilePath = flag.String("config", "", "absolute path to the config file")
	v              *viper.Viper
)

func GetConfig(configFilePath string) *viper.Viper {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		v = viper.New()
		v.SetConfigFile(configFilePath)
		err = v.ReadInConfig()
		if err != nil {
			// Log error
			panic(err)
		}
	})

	return v
}
