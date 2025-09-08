package redis

import (
	"net/url"
	"strconv"

	_ "github.com/go-viper/mapstructure/v2"
	"github.com/redis/go-redis/v9"
)

var (
	defaultConfig = redis.Options{
		PoolSize:     20,
		MinIdleConns: 5,
	}
)

type Config struct {
	Dsn         string `mapstructure:"dsn"`
	PoolSize    int    `mapstructure:"pool_size"`
	PoolMinIdle int    `mapstructure:"pool_min_idle"`
}

func NewClient(config *Config) (*redis.Client, error) {
	u, err := url.Parse(config.Dsn)
	if err != nil {
		return nil, err
	}

	db, err := strconv.Atoi(u.Path)
	if err != nil {
		return nil, err
	}

	pw, _ := u.User.Password()

	opt := &redis.Options{
		Addr:         u.Host,
		Username:     u.User.Username(),
		Password:     pw,
		DB:           db,
		PoolSize:     defaultConfig.PoolSize,
		MinIdleConns: defaultConfig.MinIdleConns,
	}
	if config.PoolSize > 0 {
		opt.PoolSize = config.PoolSize
	}
	if config.PoolMinIdle > 0 {
		opt.MinIdleConns = config.PoolMinIdle
	}

	c := redis.NewClient(opt)
	return c, nil
}
