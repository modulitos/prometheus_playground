package main

import (
	"log"

	"metrics_emitter/pkg/service"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Port  int `required:"true"`
	Debug bool
}

func main() {
	cfg := EnvConfig{}

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Panicf("could not load config, err %s", err)
	}

	sCfg := service.MetricsEmitterConfig{
		Port:        cfg.Port,
		Debug:       cfg.Debug,
		Development: true,
	}

	as := service.NewMetricsEmitterService(sCfg)
	as.Run()
}
