package scannerController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/security/jwt"
	"ticken-validator-service/utils"
)

func (controller *ScannerController) Scan(c *gin.Context) {
	validatorID := c.MustGet("jwt").(*jwt.Token).Subject

	signature := c.Param("signature")

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketScanner := controller.serviceProvider.GetTicketScanner()

	ticketScanned, err := ticketScanner.Scan(eventID, ticketID, signature, validatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	uuidVal := uuid.New()
	uuidVal.ID()

	ticketDTO := mappers.MapTicketToDTO(ticketScanned)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketDTO})
}
