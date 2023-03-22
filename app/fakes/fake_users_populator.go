package fakes

import (
	"github.com/google/uuid"
	"ticken-validator-service/config"
	"ticken-validator-service/env"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type FakeUsersPopulator struct {
	devUserInfo   config.DevUser
	reposProvider repos.IProvider
}

func NewFakeUsersPopulator(reposProvider repos.IProvider, devUserInfo config.DevUser) *FakeUsersPopulator {
	return &FakeUsersPopulator{
		devUserInfo:   devUserInfo,
		reposProvider: reposProvider,
	}
}

func (populator *FakeUsersPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	uuidDevUser, err := uuid.Parse(populator.devUserInfo.UserID)
	if err != nil {
		return err
	}

	uuidOrganization, err := uuid.Parse(populator.devUserInfo.OrganizationID)
	if err != nil {
		return err
	}

	validatorRepo := populator.reposProvider.GetValidatorRepository()

	if validatorRepo.AnyWithID(uuidDevUser) {
		return nil
	}

	devValidator := &models.Validator{
		ValidatorID:    uuidDevUser,
		Firstname:      populator.devUserInfo.Firstname,
		Lastname:       populator.devUserInfo.Lastname,
		Username:       populator.devUserInfo.Username,
		Email:          populator.devUserInfo.Email,
		OrganizationID: uuidOrganization,
	}

	if err := validatorRepo.AddValidator(devValidator); err != nil {
		return err
	}

	return nil
}
