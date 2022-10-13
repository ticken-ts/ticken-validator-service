package main

import (
	"ticken-validator-service/app"
	"ticken-validator-service/infra"
	"ticken-validator-service/utils"
)

func main() {
	tickenConfig, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	ticketValidatorService := app.New(builder, tickenConfig)
	if tickenConfig.IsDev() {
		ticketValidatorService.Populate()
	}

	ticketValidatorService.Start()
}
