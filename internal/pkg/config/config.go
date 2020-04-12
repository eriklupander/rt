package config

import "github.com/spf13/viper"

type Config struct {
	Width             int
	Height            int
	Threads           int
	DoFCamera         bool
	Samples           int
	SoftShadowSamples int
}

var Cfg *Config

func FromConfig() {
	Cfg = &Config{
		Width:             viper.GetInt("width"),
		Height:            viper.GetInt("height"),
		Threads:           viper.GetInt("threads"),
		DoFCamera:         false,
		Samples:           viper.GetInt("samples"),
		SoftShadowSamples: viper.GetInt("softshadowsamples"),
	}
}
