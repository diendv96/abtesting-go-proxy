package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type configuration struct {
	Port                     int    `mapstructure:"port"`
	Authorization            string `mapstructure:"authorization"`
	FlipEnabled              bool   `mapstructure:"flip_enabled"`
	FliptAPIEvaluate         string `mapstructure:"fliptApiEvaluate"`
	TargetSite1              string `mapstructure:"targetSite1"`
	TargetSite2              string `mapstructure:"targetSite2"`
	DefaultTargetSite        string `mapstructure:"defaultVehSite"`
	FlagKey                  string `mapstructure:"flagKey"`
	FlagCacheExpireInSeconds int    `mapstructure:"flagCacheExpireInSeconds"`
	OptimizelySDKKey         string `mapstructure:"optimizely_sdk_key"`
	FlagCacheInMb            int    `mapstructure:"flagCacheInMb"`
}

var config configuration

func init() {
	Config := viper.New()
	Config.SetConfigName("config")
	Config.SetConfigType("yml")
	Config.AddConfigPath(".")
	Config.AddConfigPath("config/")

	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	Config.AutomaticEnv()

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = Config.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Infof("Current Config: %+v", config)
}
