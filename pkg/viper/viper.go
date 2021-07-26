package viper

import (
	"github.com/spf13/viper"
)

type Viper struct {
	*viper.Viper
}

// New create viper instance by file path ex) ./config.yaml
// viper supports JSON, TOML, YAML, HCL, INI, envFile, JAVA Properties
// # docs: https://github.com/spf13/viper#reading-config-files
func NewByFilePath(filePath string) (*Viper, error) {
	v := &Viper{
		viper.New(),
	}
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Viper) Field(key string) *ViperField {
	return &ViperField{
		viper: v.Viper,
		key:   key,
	}
}
