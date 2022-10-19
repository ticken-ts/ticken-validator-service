package scannerController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-validator-service/services"
)

type ScannerController struct {
	jsonValidator   *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *ScannerController {
	controller := new(ScannerController)
	controller.jsonValidator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *ScannerController) Setup(router gin.IRouter) {
	router.POST("/events/:eventID/tickets/:ticketID/scan", controller.Scan)
}
