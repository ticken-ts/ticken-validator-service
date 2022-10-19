package main

import (
	"ticken-validator-service/app"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/infra"
)

func main() {
	tickenEnv, err := env.Load()
	if err != nil {
		panic(err)
	}

	tickenConfig, err := config.Load(tickenEnv.ConfigFilePath, tickenEnv.ConfigFileName)
	if err != nil {
		panic(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	tickenValidatorService := app.New(builder, tickenConfig)
	if tickenEnv.IsDev() {
		tickenValidatorService.Populate()
		tickenValidatorService.EmitFakeJWT()
	}

	tickenValidatorService.Start()
}
