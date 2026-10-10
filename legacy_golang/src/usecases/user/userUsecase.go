package userusecase

import (
	"fmt"

	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
	repositories "github.com/max38/golang-clean-code-architecture/src/domain/repositories"
	authentication "github.com/max38/golang-clean-code-architecture/src/shared/authentication"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	RegisterUser(userRegisterRequest *entityuser.UserRegisterRequest) (*entityuser.UserEntity, error)
	UserLogin(userLoginRequest *entityuser.UserLoginRequest) (*entityuser.UserLoginResponse, error)
	UserLogout(userRefreshTokenRequest *entityuser.UserRefreshTokenRequest) error
	UserRefreshToken(userRefreshTokenRequest *entityuser.UserRefreshTokenRequest) (*entityuser.UserLoginResponse, error)
	UserAuthentication(userToken string) (*entityuser.UserEntity, *entitymodels.UserPermissionModel, error)
}

type userUsecase struct {
	userRepository repositories.IUserRepository
}

func UserUsecase(userRepository repositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u *userUsecase) RegisterUser(userRegisterRequest *entityuser.UserRegisterRequest) (*entityuser.UserEntity, error) {
	// Hashing password
	if err := userRegisterRequest.BcryptHashing(); err != nil {
		return nil, err
	}

	var userEntity, err = u.userRepository.InsertUser(userRegisterRequest)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}

func (u *userUsecase) UserLogin(userLoginRequest *entityuser.UserLoginRequest) (*entityuser.UserLoginResponse, error) {
	var userModel, errorGetEmail = u.userRepository.GetUserByEmail(userLoginRequest.Email)
	if errorGetEmail != nil {
		return nil, errorGetEmail
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(userLoginRequest.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	var accessToken authentication.IJWTAuthentication
	var refreshToken authentication.IJWTAuthentication

	var userClaims = entityuser.UserTokenClaimsEntity{
		Id: userModel.Id,
	}
	accessToken, _ = authentication.JWTAuthentication(authentication.Access, &userClaims)
	refreshToken, _ = authentication.JWTAuthentication(authentication.Refresh, &userClaims)

	var userToken = &entityuser.UserTokenEntity{
		Access:  accessToken.SignToken(),
		Refresh: refreshToken.SignToken(),
	}

	var userLoginReponse = entityuser.UserLoginResponse{
		User:  userModel.ToEntity(),
		Token: userToken,
	}

	var errorUpsertOAuth = u.userRepository.UpsertOAuth(&userLoginReponse)
	if errorUpsertOAuth != nil {
		return nil, errorUpsertOAuth
	}
	return &userLoginReponse, nil
}

func (u *userUsecase) UserRefreshToken(userRefreshTokenRequest *entityuser.UserRefreshTokenRequest) (*entityuser.UserLoginResponse, error) {
	var claims, errorParseToken = authentication.ParseToken(userRefreshTokenRequest.RefreshToken)
	if errorParseToken != nil {
		return nil, errorParseToken
	}

	var userTokenClaim = *claims.Claims

	var userToken, errorGetTokenclaim = u.userRepository.FindOneOAuthByUserId(userTokenClaim.Id)
	if errorGetTokenclaim != nil {
		return nil, errorGetTokenclaim
	}
	if userToken == nil || userToken.Refresh != userRefreshTokenRequest.RefreshToken {
		return nil, fmt.Errorf("user token not found")
	}

	var accessToken authentication.IJWTAuthentication
	var refreshToken authentication.IJWTAuthentication

	accessToken, _ = authentication.JWTAuthentication(authentication.Access, &userTokenClaim)
	refreshToken = authentication.RepeatToken(&userTokenClaim, claims.ExpiresAt.Unix())

	var newUserToken = &entityuser.UserTokenEntity{
		Access:  accessToken.SignToken(),
		Refresh: refreshToken.SignToken(),
	}

	var userModel = userToken.User

	var userLoginReponse = entityuser.UserLoginResponse{
		User:  userModel.ToEntity(),
		Token: newUserToken,
	}

	var errorUpsertOAuth = u.userRepository.UpsertOAuth(&userLoginReponse)
	if errorUpsertOAuth != nil {
		return nil, errorUpsertOAuth
	}
	return &userLoginReponse, nil
}

func (u *userUsecase) UserLogout(userRefreshTokenRequest *entityuser.UserRefreshTokenRequest) error {
	var claims, errorParseToken = authentication.ParseToken(userRefreshTokenRequest.RefreshToken)
	if errorParseToken != nil {
		return errorParseToken
	}

	var userTokenClaim = *claims.Claims

	var userToken, errorGetTokenclaim = u.userRepository.FindOneOAuthByUserId(userTokenClaim.Id)
	if errorGetTokenclaim != nil {
		return errorGetTokenclaim
	}
	if userToken == nil || userToken.Refresh != userRefreshTokenRequest.RefreshToken {
		return fmt.Errorf("user token not found")
	}

	var errorDeleteOAuth = u.userRepository.DeleteOAuthByUserId(userTokenClaim.Id)
	if errorDeleteOAuth != nil {
		return errorDeleteOAuth
	}
	return nil
}

func (u *userUsecase) UserAuthentication(userToken string) (*entityuser.UserEntity, *entitymodels.UserPermissionModel, error) {
	var claims, errorParseToken = authentication.ParseToken(userToken)
	if errorParseToken != nil {
		return nil, nil, errorParseToken
	}

	var userTokenClaim = *claims.Claims

	var userTokenModel, errorGetTokenclaim = u.userRepository.FindOAuthByUserIdAndAccessToken(userTokenClaim.Id, userToken)

	if errorGetTokenclaim != nil {
		return nil, nil, errorGetTokenclaim
	}
	if userTokenModel == nil {
		return nil, nil, fmt.Errorf("user token not found")
	}

	var userPermission, _ = u.userRepository.GetUserPermissionByUserId(userTokenModel.UserId)

	var userEntity = userTokenModel.User

	return userEntity.ToEntity(), userPermission, nil
}

func (u *userUsecase) GetUserProfile(userId uint) (*entityuser.UserEntity, error) {
	var userEntity, errorGetUser = u.userRepository.GetUserByUserId(userId)
	if errorGetUser != nil {
		return nil, errorGetUser
	}

	return userEntity.ToEntity(), nil
}
