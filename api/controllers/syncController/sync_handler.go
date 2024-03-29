package syncController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-validator-service/utils"
)

func (controller *SyncController) Sync(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketSyncer := controller.serviceProvider.GetTicketSyncer()

	// TODO -> user from organization
	if err := ticketSyncer.Sync(eventID, uuid.Nil); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, utils.HttpResponse{Message: "ticket syncer process started"})
}
