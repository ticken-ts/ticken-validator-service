package validatorsController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-validator-service/services"
)

type ValidatorsController struct {
	jsonValidator   *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *ValidatorsController {
	controller := new(ValidatorsController)
	controller.jsonValidator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *ValidatorsController) Setup(router gin.IRouter) {
	router.POST("/organizations/:organizationID/validators", controller.RegisterValidator)
}
