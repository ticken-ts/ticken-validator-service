package main

import (
	"ticken-validator-service/app"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/infra"
	"ticken-validator-service/log"
)

func main() {
	tickenEnv, err := env.Load()
	if err != nil {
		panic(err)
	}

	log.InitGlobalLogger()

	tickenConfig, err := config.Load(tickenEnv.ConfigFilePath, tickenEnv.ConfigFileName)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}

	tickenValidatorService := app.New(builder, tickenConfig)
	tickenValidatorService.Populate()

	if tickenEnv.IsDev() {
		tickenValidatorService.EmitFakeJWT()
	}

	tickenValidatorService.Start()
}
