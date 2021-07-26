package viper

import (
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type ViperField struct {
	viper *viper.Viper
	key   string
}

func (v *ViperField) Default(value interface{}) *ViperField {
	v.viper.SetDefault(v.key, value)
	return v
}

func (v *ViperField) Bool() bool {
	return cast.ToBool(v.viper.Get(v.key))
}

func (v *ViperField) Int() int {
	return cast.ToInt(v.viper.Get(v.key))
}

func (v *ViperField) Int64() int64 {
	return cast.ToInt64(v.viper.Get(v.key))
}

func (v *ViperField) Uint64() uint64 {
	return cast.ToUint64(v.viper.Get(v.key))
}

func (v *ViperField) String() string {
	return cast.ToString(v.viper.Get(v.key))
}

func (v *ViperField) StringSlice() []string {
	return cast.ToStringSlice(v.viper.Get(v.key))
}

func (v *ViperField) StringMapString() map[string]string {
	return cast.ToStringMapString(v.viper.Get(v.key))
}

func (v *ViperField) Duration() time.Duration {
	return cast.ToDuration(v.viper.Get(v.key))
}
