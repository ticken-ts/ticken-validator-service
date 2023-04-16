package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
	"ticken-validator-service/api"
	"ticken-validator-service/api/controllers/healthController"
	"ticken-validator-service/api/controllers/scannerController"
	"ticken-validator-service/api/controllers/validatorsController"
	"ticken-validator-service/api/middlewares"
	"ticken-validator-service/app/fakes"
	"ticken-validator-service/async"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/infra"
	"ticken-validator-service/log"
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
	jwtVerifier     jwt.Verifier
	subscriber      *async.Subscriber

	// populators are intended to populate
	// useful data. It can be testdata or
	// data that should be present on the db
	// before the service is available
	populators []Populator
}

func New(infraBuilder *infra.Builder, tickenConfig *config.Config) *TickenValidatorApp {
	log.TickenLogger.Info().Msg(fmt.Sprintf("initializing service: %s", color.BlueString(tickenConfig.Server.ClientID)))

	ticketValidatorApp := new(TickenValidatorApp)

	/******************************** infra builds ********************************/
	engine := infraBuilder.BuildEngine()
	pvtbcCaller := infraBuilder.BuildPvtbcCaller()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	pubbcCaller := infraBuilder.BuildPubbcCaller(env.TickenEnv.TickenWalletKey)
	busSubscriber := infraBuilder.BuildBusSubscriber(env.TickenEnv.BusConnString)
	/*******************************************************************************/

	/********************************** providers **********************************/
	repoProvider, err := repos.NewProvider(
		db,
		&tickenConfig.Database,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}

	serviceProvider, err := services.NewProvider(
		repoProvider,
		infraBuilder,
		pvtbcCaller,
		pubbcCaller,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}
	/*******************************************************************************/

	/********************************* subscriber **********************************/
	subscriber, err := async.NewSubscriber(busSubscriber, serviceProvider)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}
	/*******************************************************************************/

	ticketValidatorApp.engine = engine
	ticketValidatorApp.config = tickenConfig
	ticketValidatorApp.subscriber = subscriber
	ticketValidatorApp.jwtVerifier = jwtVerifier
	ticketValidatorApp.repoProvider = repoProvider
	ticketValidatorApp.serviceProvider = serviceProvider

	ticketValidatorApp.loadMiddlewares(engine)
	ticketValidatorApp.loadControllers(engine)

	/********************************* populators **********************************/
	ticketValidatorApp.populators = []Populator{
		fakes.NewFakeUsersPopulator(repoProvider, tickenConfig.Dev.User),
		fakes.NewFakeTicketsPopulator(repoProvider, tickenConfig.Dev.User),
	}
	/*******************************************************************************/

	return ticketValidatorApp
}

func (tickenValidatorApp *TickenValidatorApp) Start() {
	if err := tickenValidatorApp.subscriber.Start(); err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to subscriber : %s", err.Error()))
	}

	url := tickenValidatorApp.config.Server.GetServerURL()
	if err := tickenValidatorApp.engine.Run(url); err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to start server: %s", err.Error()))
	}
}

func (tickenValidatorApp *TickenValidatorApp) Populate() {
	for _, populator := range tickenValidatorApp.populators {
		err := populator.Populate()
		if err != nil {
			panic(err)
		}
	}
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

func (tickenValidatorApp *TickenValidatorApp) loadControllers(apiRouter gin.IRouter) {
	apiRouterGroup := apiRouter.Group(tickenValidatorApp.config.Server.APIPrefix)

	var appControllers = []api.Controller{
		healthController.New(tickenValidatorApp.serviceProvider),
		scannerController.New(tickenValidatorApp.serviceProvider),
		validatorsController.New(tickenValidatorApp.serviceProvider),
	}

	for _, controller := range appControllers {
		controller.Setup(apiRouterGroup)
	}
}

func (tickenValidatorApp *TickenValidatorApp) loadMiddlewares(apiRouter gin.IRouter) {
	var appMiddlewares = []api.Middleware{
		middlewares.NewErrorMiddleware(),
		middlewares.NewAuthMiddleware(tickenValidatorApp.jwtVerifier, tickenValidatorApp.config.Server.APIPrefix),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}
}
