package app

import (
	"github.com/gin-gonic/gin"
	"ticken-validator-service/api"
	"ticken-validator-service/api/controllers/scannerController"
	"ticken-validator-service/api/middlewares"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos"
	"ticken-validator-service/services"
)

type TickenValidatorApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
}

func New(builder *infra.Builder, tickenConfig *config.Config) *TickenValidatorApp {
	ticketValidatorApp := new(TickenValidatorApp)

	engine := builder.BuildEngine()
	pvtbcCaller := builder.BuildPvtbcCaller()
	db := builder.BuildDb(env.TickenEnv.ConnString)

	// this provider is going to provider all repositories
	// to the services
	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(repoProvider, pvtbcCaller)
	if err != nil {
		panic(err)
	}

	ticketValidatorApp.engine = engine
	ticketValidatorApp.repoProvider = repoProvider
	ticketValidatorApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	var controllers = []api.Controller{
		scannerController.New(serviceProvider),
	}

	for _, controller := range controllers {
		controller.Setup(engine)
	}

	return ticketValidatorApp
}

func (tickenValidatorApp *TickenValidatorApp) Start() {
	url := tickenValidatorApp.config.Server.GetServerURL()
	err := tickenValidatorApp.engine.Run(url)
	if err != nil {
		panic(err)
	}
}

func (tickenValidatorApp *TickenValidatorApp) Populate() {
}
