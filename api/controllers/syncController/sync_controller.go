package syncController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-validator-service/services"
)

type SyncController struct {
	jsonValidator   *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *SyncController {
	controller := new(SyncController)
	controller.jsonValidator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *SyncController) Setup(router gin.IRouter) {
	router.POST("/events/:eventID/sync", controller.Sync)
}
