package postgresuserrepository

import (
	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
	repositories "github.com/max38/golang-clean-code-architecture/src/domain/repositories"
	postgreshandler "github.com/max38/golang-clean-code-architecture/src/infrastructure/database/postgres"
)

type userRepository struct {
	dbHandler postgreshandler.IPostgresHandler
}

func UserRepository() repositories.IUserRepository {
	var dbHandler = postgreshandler.PostgresHandler()
	return &userRepository{
		dbHandler: dbHandler,
	}
}

func (u *userRepository) InsertUser(userRegisterRequest *entityuser.UserRegisterRequest) (*entityuser.UserEntity, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userModel = entitymodels.UserModel{
		UserEntity: entityuser.UserEntity{
			Email:     userRegisterRequest.Email,
			FirstName: userRegisterRequest.FirstName,
			LastName:  userRegisterRequest.LastName,
		},
		Password: userRegisterRequest.Password,
	}

	var dbResult = dbConnector.Create(&userModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	var userEntity = *userModel.ToEntity()

	return &userEntity, nil
}

func (u *userRepository) GetUserByEmail(email string) (*entitymodels.UserModel, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userModel = entitymodels.UserModel{}
	var dbResult = dbConnector.Where("email = ?", email).First(&userModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &userModel, nil
}

func (u *userRepository) GetUserByUserId(userId uint) (*entitymodels.UserModel, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userModel = entitymodels.UserModel{}
	var dbResult = dbConnector.Where("id = ?", userId).First(&userModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &userModel, nil
}

func (u *userRepository) InsertOAuth(userLoginReponse *entityuser.UserLoginResponse) error {
	var dbConnector = u.dbHandler.GetConnector()

	var userTokenModel = entitymodels.UserTokenModel{
		UserId:  userLoginReponse.User.Id,
		Access:  userLoginReponse.Token.Access,
		Refresh: userLoginReponse.Token.Refresh,
	}

	var dbResult = dbConnector.Create(&userTokenModel)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}

func (u *userRepository) UpsertOAuth(userLoginReponse *entityuser.UserLoginResponse) error {
	var dbConnector = u.dbHandler.GetConnector()

	var userTokenModel = entitymodels.UserTokenModel{
		UserId:  userLoginReponse.User.Id,
		Access:  userLoginReponse.Token.Access,
		Refresh: userLoginReponse.Token.Refresh,
	}

	var dbResult = dbConnector.Where("user_id = ?", userLoginReponse.User.Id).Updates(&userTokenModel)
	if dbResult.RowsAffected == 0 {
		dbResult = dbConnector.Create(&userTokenModel)
	}
	if dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}

func (u *userRepository) FindOneOAuthByUserId(userId uint) (*entitymodels.UserTokenModel, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userTokenModel = entitymodels.UserTokenModel{}
	var dbResult = dbConnector.Joins("User").Where("user_id = ?", userId).First(&userTokenModel)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &userTokenModel, nil
}

func (u *userRepository) DeleteOAuthByUserId(userId uint) error {
	var dbConnector = u.dbHandler.GetConnector()

	var dbResult = dbConnector.Where("user_id = ?", userId).Delete(&entitymodels.UserTokenModel{})
	if dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}

func (u *userRepository) FindOAuthByUserIdAndAccessToken(userId uint, accessToken string) (*entitymodels.UserTokenModel, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userTokenModel = entitymodels.UserTokenModel{}
	var dbResult = dbConnector.Joins("User").Where("user_id = ? AND access = ?", userId, accessToken).First(&userTokenModel)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &userTokenModel, nil
}

func (u *userRepository) GetUserPermissionByUserId(userId uint) (*entitymodels.UserPermissionModel, error) {
	var dbConnector = u.dbHandler.GetConnector()

	var userPermissionModel = entitymodels.UserPermissionModel{}
	var dbResult = dbConnector.Where("user_id = ?", userId).First(&userPermissionModel)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &userPermissionModel, nil
}
