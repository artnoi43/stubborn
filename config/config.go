package config

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/artnoi43/stubborn/cmd"
	"github.com/artnoi43/stubborn/data/repository"
	"github.com/artnoi43/stubborn/domain/usecase/dnsserver"
	"github.com/artnoi43/stubborn/domain/usecase/handler"
)

type Config struct {
	ServerConfig  dnsserver.Config  `mapstructure:"server"`
	CacherConfig  repository.Config `mapstructure:"cacher"`
	HandlerConfig handler.Config    `mapstructure:"handler"`
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

func InitConfig(confLocation string) (conf *Config, err error) {
	dir, fullFilename := path.Split(confLocation)
	fileAndExt := strings.Split(fullFilename, ".")
	if len(fileAndExt) < 2 {
		log.Fatalln("bad config file location:", confLocation)
	}
	filename := fileAndExt[0]
	fileExtension := fileAndExt[1]
	// Defaults
	viper.SetDefault("handler.hosts_file", cmd.TableLocation)
	viper.SetDefault("handler.all_types", true)
	viper.SetDefault("handler.dot.outbound", "DOT")
	viper.SetDefault("handler.dot.upstream_timeout", 10)
	viper.SetDefault("handler.dot.upstream_ip", "1.1.1.1")
	viper.SetDefault("handler.dot.upstream_port", "853")
	viper.SetDefault("server.address", "127.0.0.1:5300")
	viper.SetDefault("server.protocol", "udp")
	viper.SetDefault("cacher.expiration", 300)
	viper.SetDefault("cacher.cleanup", 600)

	err = loadConf(dir, filename, fileExtension)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading config file")
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing config file")
	}
	if !conf.HandlerConfig.Outbound.IsValid() {
		// For readability
		msg := "bad outbound in config file\noutbound: \"%s\"\nconfig location: %s\n"
		return nil, fmt.Errorf(msg, conf.HandlerConfig.Outbound, confLocation)
	}
	return conf, nil
}

func loadConf(dir string, file string, ext string) error {
	// Default config file dir is /etc/stubborn/config.yaml
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
