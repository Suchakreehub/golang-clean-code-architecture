package config

import (
	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
)

var EntitiyModels = []entities.ICRUDDataModel{
	// PostgreSQL
	&entitymodels.UserModel{},
	&entitymodels.UserTokenModel{},
	&entitymodels.UserPermissionModel{},

	// MongoDB
}
