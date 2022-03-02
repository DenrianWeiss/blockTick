package config

import "github.com/spf13/viper"

var Config Conf

type Conf struct {
	DefaultRpc      string          `json:"default_rpc"`
	EnablePow       bool            `json:"enable_pow"`
	PowDifficulty   int             `json:"pow_difficulty"`
	EnableWhitelist bool            `json:"enable_whitelist"`
	Whitelist       map[string]bool `json:"whitelist"`
	EnableSecret    bool            `json:"enable_secret"`
	Secret          string          `json:"secret"`
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("No config")
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
