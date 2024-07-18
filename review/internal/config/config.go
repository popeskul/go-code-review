package config

import (
	"coupon_service/internal/api"
	"log"

	"github.com/brumhard/alligotor"
)

type Config struct {
	API api.Config `yaml:"api"`
}

func New() Config {
	cfg := Config{}
	if err := alligotor.Get(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
