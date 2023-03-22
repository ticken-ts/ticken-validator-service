package validatorsController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/utils"
)

type registerValidatorPayload struct {
	JWT string `json:"jwt"`
}

func (controller *ValidatorsController) RegisterValidator(c *gin.Context) {
	var payload registerValidatorPayload

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	validatorManager := controller.serviceProvider.GetValidatorManager()

	validator, err := validatorManager.RegisterValidator(organizationID, payload.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HttpResponse{Message: err.Error()})
		return
	}

	c.JSON(
		http.StatusCreated,
		utils.HttpResponse{
			Message: "Validator created successfully",
			Data:    mappers.MapValidatorToDTO(validator),
		},
	)
}
