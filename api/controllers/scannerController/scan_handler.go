package scannerController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/security/jwt"
	"ticken-validator-service/utils"
)

type scanTicketPayload struct {
	RSignatureField string `json:"r"`
	SSignatureField string `json:"s"`
}

func (controller *ScannerController) Scan(c *gin.Context) {
	validatorID := c.MustGet("jwt").(*jwt.Token).Subject

	var payload scanTicketPayload

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

	if err = c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	ticketScanner := controller.serviceProvider.GetTicketScanner()

	ticketScanned, err := ticketScanner.Scan(
		eventID,
		ticketID,
		validatorID,
		payload.RSignatureField,
		payload.SSignatureField,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketDTO := mappers.MapTicketToDTO(ticketScanned)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketDTO})
}
