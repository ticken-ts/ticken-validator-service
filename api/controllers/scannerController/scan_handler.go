package scannerController

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/security/jwt"
	"ticken-validator-service/utils"
)

func (controller *ScannerController) Scan(c *gin.Context) {
	validatorID := c.MustGet("jwt").(*jwt.Token).Subject

	b64signature := c.Query("signature")

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

	signature, err := base64.URLEncoding.DecodeString(b64signature)
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

	ticketDTO := mappers.MapTicketToDTO(ticketScanned)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketDTO})
}
