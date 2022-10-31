package dep_container

import (
	"crypto-to-fiat-converter/intenral/config"
	"crypto-to-fiat-converter/pkg/logger"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

const loggerDefName = "logger"

// RegisterLogger registers Logger dependency.
func RegisterLogger(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: loggerDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)

			return logger.Init(cfg.LogLevel), nil
		},
		Close: func(obj interface{}) error {
			obj.(*zap.Logger).Sync() // nolint: errcheck
			return nil
		},
	})
}

// InitLogger init Logger.
func (c Container) InitLogger() {
	c.container.Get(loggerDefName)
}
