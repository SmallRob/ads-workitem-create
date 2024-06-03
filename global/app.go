package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"adswork/config"
)

//
type Application struct {
	ConfigViper *viper.Viper
	// 全局配置
	Config config.Configuration
	Log    *zap.Logger
}

// App 实例化对象
var App = new(Application)
