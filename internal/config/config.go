package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Conf struct {
	Port              string        `mapstructure:"PORT"`
	RedisAddr         string        `mapstructure:"REDIS_ADDR"`
	RedisPassword     string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int           `mapstructure:"REDIS_DB"`
	DefaultLimitIP    int           `mapstructure:"DEFAULT_LIMIT_IP"`
	DefaultLimitToken int           `mapstructure:"DEFAULT_LIMIT_TOKEN"`
	BlockDurationMin  int           `mapstructure:"BLOCK_DURATION_MINUTES"`
	BlockDuration     time.Duration

func LoadConfig(path string) (*Conf, error) {
	var cfg Conf

	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("DEFAULT_LIMIT_IP", 10)
	viper.SetDefault("DEFAULT_LIMIT_TOKEN", 100)
	viper.SetDefault("BLOCK_DURATION_MINUTES", 5)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Buscando variáveis da memória...")
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.BlockDuration = time.Duration(cfg.BlockDurationMin) * time.Minute

	return &cfg, nil
}