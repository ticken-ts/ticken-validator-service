package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
	"ticken-validator-service/api"
	"ticken-validator-service/api/controllers/healthController"
	"ticken-validator-service/api/controllers/scannerController"
	"ticken-validator-service/api/middlewares"
	"ticken-validator-service/async"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos"
	"ticken-validator-service/security/jwt"
	"ticken-validator-service/services"
	"ticken-validator-service/utils"
)

type TickenValidatorApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
}

func New(infraBuilder *infra.Builder, tickenConfig *config.Config) *TickenValidatorApp {
	ticketValidatorApp := new(TickenValidatorApp)

	engine := infraBuilder.BuildEngine()
	pvtbcCaller := infraBuilder.BuildPvtbcCaller()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	pubbcCaller := infraBuilder.BuildPubbcCaller(env.TickenEnv.TickenWalletKey)
	busSubscriber := infraBuilder.BuildBusSubscriber(env.TickenEnv.BusConnString)

	// this provider is going to provider all repositories
	// to the services
	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(repoProvider, pvtbcCaller, pubbcCaller)
	if err != nil {
		panic(err)
	}

	subscriber, err := async.NewSubscriber(busSubscriber, serviceProvider)
	if err != nil {
		panic(err)
	}

	err = subscriber.Start()
	if err != nil {
		panic(err)
	}

	ticketValidatorApp.engine = engine
	ticketValidatorApp.config = tickenConfig
	ticketValidatorApp.repoProvider = repoProvider
	ticketValidatorApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, jwtVerifier, tickenConfig.Server.APIPrefix),
	}

	var appControllers = []api.Controller{
		healthController.New(serviceProvider),
		scannerController.New(serviceProvider),
	}

	apiRouter := engine.Group(tickenConfig.Server.APIPrefix)

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}

	for _, controller := range appControllers {
		controller.Setup(apiRouter)
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
	rsaPrivKey, err := utils.LoadRSA(
		tickenValidatorApp.config.Dev.JWTPrivateKey, tickenValidatorApp.config.Dev.JWTPublicKey)
	if err != nil {
		panic(err)
	}

	fakeJWT := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &jwt.Claims{
		Subject:           tickenValidatorApp.config.Dev.User.UserID,
		Email:             tickenValidatorApp.config.Dev.User.Email,
		PreferredUsername: tickenValidatorApp.config.Dev.User.Username,
	})

	signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

	if err != nil {
		panic(fmt.Errorf("error generation fake Token: %s", err.Error()))
	}

	fmt.Printf("DEV Token: %s \n", signedJWT)
}
