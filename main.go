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

	tickenConfig, err := config.Load(".")
	if err != nil {
		panic(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	ticketValidatorService := app.New(builder, tickenConfig)
	if tickenEnv.IsDev() {
		ticketValidatorService.Populate()
	}

	ticketValidatorService.Start()
}
