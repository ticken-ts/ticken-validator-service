package scannerController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/utils"
)

func (controller *ScannerController) Scan(c *gin.Context) {
	eventID, ticketID := c.Param("eventID"), c.Param("ticketID")

	ticketScanner := controller.serviceProvider.GetTicketScanner()

	ticketScanned, err := ticketScanner.Scan(eventID, ticketID, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketScannedDto := mappers.MapChainTicketToDTO(ticketScanned)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketScannedDto})
}
