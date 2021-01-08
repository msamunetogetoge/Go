package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	Host      string
	Key       string
	LogFile   string
	DbName    string
	SQLDriver string
	Port      int
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		Host:    cfg.Section("rakuten").Key("x-rapidapi-host").String(),
		Key:     cfg.Section("rakuten").Key("x-rapidapi-key").String(),
		LogFile: cfg.Section("weather").Key("log_file").String(),
		Port:    cfg.Section("web").Key("port").MustInt(),
	}
}
