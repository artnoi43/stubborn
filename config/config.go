package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/artnoi43/stubborn/lib/cacher"
	"github.com/artnoi43/stubborn/lib/dnsserver"
	"github.com/artnoi43/stubborn/lib/handler"
	"github.com/artnoi43/stubborn/lib/rediswrapper"
)

type Config struct {
	ServerConfig  dnsserver.Config    `mapstructure:"server"`
	CacherConfig  cacher.Config       `mapstructure:"cacher"`
	RedisConfig   rediswrapper.Config `mapstructure:"redis"`
	HandlerConfig handler.Config      `mapstructure:"handler"`
}

type Location struct {
	Dir  string
	Name string
	Ext  string
}

func ParsePath(rawPath string) *Location {
	dir, configFile := filepath.Split(rawPath)
	name := strings.Split(configFile, ".")[0] // remove ext from filename
	ext := filepath.Ext(configFile)[1:]       // remove dot
	return &Location{
		Dir:  dir,
		Name: name,
		Ext:  ext,
	}
}

func InitConfig(dir string, file string, ext string) (conf *Config, err error) {
	// Defaults
	viper.SetDefault("handler.hosts_file", "./config/table.json")
	viper.SetDefault("handler.all_types", true)
	viper.SetDefault("server.address", "127.0.0.1:5300")
	viper.SetDefault("server.protocol", "udp")
	viper.SetDefault("cacher.expiration", 300)
	viper.SetDefault("cacher.cleanup_interval", 600)
	viper.SetDefault("redis.address", "127.0.0.1:6379")
	viper.SetDefault("redis.username", "")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 1)

	err = loadConf(dir, file, ext)
	if err != nil {
		return nil, err
	}
	conf, err = unmarshal()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func loadConf(dir string, file string, ext string) error {
	// Default config file dir is $HOME/config/fngobot
	// From CLI: -c <path>
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType(ext)

	// Parse config from both env and file
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		// Config file not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// WriteConfig() just won't create new file if doesn't exist
			viper.SafeWriteConfig()
		} else {
			return err
		}
	}
	return nil
}

func unmarshal() (conf *Config, err error) {
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
