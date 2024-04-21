package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	BotToken  string             `mapstructure:"botToken"`
	Logger    Logger             `mapstructure:"logger"`
	Admins    map[int64]AdminFlags `mapstructure:"admins"`
	Servers   map[string]Server `mapstructure:"servers"`
	Databases map[string]DB     `mapstructure:"databases"`
}

type Logger struct {
	BotDebug    bool   `mapstructure:"botDebug"`
	Development bool   `mapstructure:"development"`
	Encoding    string `mapstructure:"encoding"`
}

type Server struct {
	IP        string            `mapstructure:"ip"`
	Port      int               `mapstructure:"port"`
	Databases map[string]DBName `mapstructure:"databases"`
	Rcon 	  Rcon     `mapstructure:"rcon"`
}

type DBName struct {
	DB    string `mapstructure:"db"`
	ServerID    string `mapstructure:"server_id"`
	Table string `mapstructure:"table"`
}

type Rcon struct {
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type AdminFlags struct {
	Flags string `mapstructure:"flags"`
}

func NewConfig(configName, configPath string) (*Config, error) {
	var cfg Config

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return &cfg, nil
}
