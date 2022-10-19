package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
		middlewares.NewAuthMiddleware(serviceProvider, &tickenConfig.Server),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	var appControllers = []api.Controller{
		scannerController.New(serviceProvider),
	}

	for _, controller := range appControllers {
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

func (tickenValidatorApp *TickenValidatorApp) EmitFakeJWT() {
	fakeJWT := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":   "290c641a-55a1-40f5-acc3-d4ebe3626fdd",
		"email": "user@ticken.com",
	})

	signedJWT, err := fakeJWT.SigningString()
	if err != nil {
		panic(fmt.Errorf("error generation fake JWT: %s", err.Error()))
	}

	fmt.Printf("DEV JWT: %s \n", signedJWT)
}
