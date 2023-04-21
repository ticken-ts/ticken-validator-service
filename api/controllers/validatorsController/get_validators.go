package validatorsController

import (
	"net/http"
	"ticken-validator-service/api/mappers"
	"ticken-validator-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (controller *ValidatorsController) GetValidators(c *gin.Context) {
	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	validatorManager := controller.serviceProvider.GetValidatorManager()

	validators, err := validatorManager.GetValidators(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HttpResponse{Message: err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		utils.HttpResponse{
			Message: "Validators retrieved successfully",
			Data:    mappers.MapValidatorsToDTOs(validators),
		},
	)
}
