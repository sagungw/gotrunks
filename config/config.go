package config

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	ENV_LOCAL      = "local"
	ENV_PRODUCTION = "production"
	ENV_STAGING    = "staging"
)

type EnvAware interface {
	GetEnv() string
}

type Application struct {
	Env     string `mapstructure:"env"`
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type PostLoader[T EnvAware] func(*T) error

type LoadConfiguration[T EnvAware] struct {
	ConfigName  string
	ConfigPaths []string
	EnvPrefix   string
	PostLoaders []PostLoader[T]

	configName string
	configType string
	once       sync.Once
}

func (l *LoadConfiguration[T]) SeparateConfigNameAndType() (string, string) {
	l.once.Do(func() {
		ss := strings.Split(l.ConfigName, ".")
		if len(ss) < 2 {
			return
		}

		l.configName = ss[0]
		l.configType = ss[1]
	})

	return l.configName, l.configType
}

func Load[T EnvAware](ctx context.Context, loadConfig *LoadConfiguration[T]) (*T, error) {
	if loadConfig == nil {
		return nil, errors.New("load configuration cannot be nil")
	}

	var configuration T
	configName, configType := loadConfig.SeparateConfigNameAndType()
	if configName == "" || configType == "" {
		err := errors.New("config file name is invalid")
		return nil, err
	}

	viper := viper.New()
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(loadConfig.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, path := range loadConfig.ConfigPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	for _, f := range loadConfig.PostLoaders {
		if err := f(&configuration); err != nil {
			return nil, err
		}
	}

	return &configuration, nil
}
