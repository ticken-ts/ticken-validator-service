package app

import (
	"ticken-validator-service/api"
	"ticken-validator-service/api/controllers/scannerController"
	"ticken-validator-service/api/middlewares"
	"ticken-validator-service/infra"
	"ticken-validator-service/services"
	"ticken-validator-service/utils"
)

type TickenValidatorApp struct {
	router          infra.Router
	serviceProvider services.Provider
}

func New(router infra.Router, db infra.Db, tickenConfig *utils.TickenConfig) *TickenValidatorApp {
	ticketValidatorApp := new(TickenValidatorApp)

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, _ := services.NewProvider(db, tickenConfig)

	ticketValidatorApp.router = router
	ticketValidatorApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(router)
	}

	var controllers = []api.Controller{
		scannerController.New(serviceProvider),
	}

	for _, controller := range controllers {
		controller.Setup(router)
	}

	return ticketValidatorApp
}

func (tickenTicketApp *TickenValidatorApp) Start() {
	err := tickenTicketApp.router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func (tickenTicketApp *TickenValidatorApp) Populate() {
}
