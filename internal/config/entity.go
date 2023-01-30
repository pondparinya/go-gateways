package config

type (
	Configs struct {
		Port string `mapstructure:"Port"`
	}
)

var APP = &Configs{}
