package scannerController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/utils"
)

func (controller *ScannerController) Scan(c *gin.Context) {
	eventID, ticketID := c.Param("eventID"), c.Param("ticketID")
	owner := c.MustGet("jwt").(*oidc.IDToken).Subject

	ticketScanner := controller.serviceProvider.GetTicketScanner()

	ticketScanned, err := ticketScanner.Scan(eventID, ticketID, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketDTO := mappers.MapTicketToDTO(ticketScanned)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketDTO})
}
