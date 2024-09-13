package services

import "github.com/spf13/viper"

func GetConfig(configFilePath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(configFilePath)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}
