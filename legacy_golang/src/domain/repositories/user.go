package repositories

import (
	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*entitymodels.UserModel, error)
	GetUserByUserId(userId uint) (*entitymodels.UserModel, error)
	InsertUser(userRegisterRequest *entityuser.UserRegisterRequest) (*entityuser.UserEntity, error)
	InsertOAuth(userLoginReponse *entityuser.UserLoginResponse) error
	UpsertOAuth(userLoginReponse *entityuser.UserLoginResponse) error
	FindOneOAuthByUserId(userId uint) (*entitymodels.UserTokenModel, error)
	FindOAuthByUserIdAndAccessToken(userId uint, accessToken string) (*entitymodels.UserTokenModel, error)
	DeleteOAuthByUserId(userId uint) error
	GetUserPermissionByUserId(userId uint) (*entitymodels.UserPermissionModel, error)
}
