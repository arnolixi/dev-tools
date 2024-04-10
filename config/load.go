package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(Cfg any, configPath string) *viper.Viper {
	v := viper.New()
	if configPath == "" {
		panic("not Found config file")
	}
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		zap.L().Error("Load Config file failed.", zap.Error(err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("config file changed.", zap.Any("eventName", e.Name), zap.Any("eventStr", e.String()))
		if err = v.Unmarshal(Cfg); err != nil {
			zap.L().Error("OnConfigChange Unmarshal failed ", zap.Error(err))
		}
	})
	if err = v.Unmarshal(Cfg); err != nil {
		zap.L().Error("Load Unmarshal failed ", zap.Error(err))
	}
	return v
}
