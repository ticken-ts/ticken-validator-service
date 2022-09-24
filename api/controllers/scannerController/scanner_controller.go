package scannerController

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"ticken-validator-service/infra"
	"ticken-validator-service/services"
)

type ScannerController struct {
	jsonValidator   *validator.Validate
	serviceProvider services.Provider
}

// TODO -> test only until user management is complete
var owner = uuid.New().String()

func New(serviceProvider services.Provider) *ScannerController {
	controller := new(ScannerController)
	controller.jsonValidator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *ScannerController) Setup(router infra.Router) {
	router.POST("/events/:eventID/tickets/:ticketID/scan", controller.Scan)
}
