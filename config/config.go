package config

import (
	"echoapptpl/util"
	"fmt"
	"github.com/spf13/viper"
)

var Cfg Config

type Config map[string]interface{}

func Init() error {
	// 解析到Config
	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}
	fmt.Println("cfg:", util.JsonPretty(&Cfg))
	return nil
}
