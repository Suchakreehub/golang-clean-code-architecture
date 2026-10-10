package gofiberuserhandler

import (
	"github.com/gofiber/fiber/v2"

	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
	gofiberentities "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/entities"
	userusecase "github.com/max38/golang-clean-code-architecture/src/usecases/user"
)

type userHandlersErrorCode string

const (
	registerUserError userHandlersErrorCode = "users-001"
	loginUserError    userHandlersErrorCode = "users-002"
	refreshTokenError userHandlersErrorCode = "users-003"
	logoutUserError   userHandlersErrorCode = "users-004"
)

type IUserHandler interface {
	RegisterUser(c *fiber.Ctx) error
	UserLogin(c *fiber.Ctx) error
	UserLogout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type userHandler struct {
	userUsecase userusecase.IUserUsecase
}

func UserHandler(userUsecase userusecase.IUserUsecase) IUserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

// @Summary Register User
// @Description Register User
// @Tags User
// @Accept application/json
// @Produce json
// @Success 201 {object} entityuser.UserEntity
// @Failure 400 {object} gofiberentities.ErrorResponseType
// @Router /api/v1/register [post]
// @Param request body entityuser.UserRegisterRequest true "query params"
func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	// Parse request
	var userRegisterRequest = new(entityuser.UserRegisterRequest) // Pointer
	if err := c.BodyParser(userRegisterRequest); err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(registerUserError),
			err.Error(),
		).Response()
	}

	// Email validation
	if !userRegisterRequest.IsEmail() {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(registerUserError),
			"email pattern is invalid",
		).Response()
	}

	// Pass thorugh Use Case
	var userEntity, err = h.userUsecase.RegisterUser(userRegisterRequest)
	if err != nil {
		// Email has been used
		// insert user error
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(registerUserError),
			err.Error(),
		).Response()
	}

	// Return response
	return gofiberentities.NewResponse(c).Success(fiber.StatusCreated, userEntity).Response()
}

// UserLogin handles user login.
// @Summary User Login
// @Description User Login
// @Tags User
// @Accept application/json
// @Produce json
// @Param request body entityuser.UserLoginRequest true "query params"
// @Success 200 {object} entityuser.UserLoginResponse
// @Failure 400 {object} gofiberentities.ErrorResponseType
// @Router /api/v1/login [post]
func (h *userHandler) UserLogin(c *fiber.Ctx) error {
	// Parse request
	var userLoginRequest = new(entityuser.UserLoginRequest) // Pointer
	if err := c.BodyParser(userLoginRequest); err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(loginUserError),
			err.Error(),
		).Response()
	}

	// Pass thorugh Use Case
	var userLoginReponse, err = h.userUsecase.UserLogin(userLoginRequest)
	if err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(loginUserError),
			err.Error(),
		).Response()
	}

	// Return response
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, userLoginReponse).Response()
}

// UserLogout handles user logout.
// @Summary User Logout
// @Description User Logout
// @Tags User
// @Accept application/json
// @Produce json
// @Param request body entityuser.UserRefreshTokenRequest true "query params"
// @Success 200
// @Failure 400 {object} gofiberentities.ErrorResponseType
// @Router /api/v1/logout [post]
func (h *userHandler) UserLogout(c *fiber.Ctx) error {
	var userRefreshTokenRequest = new(entityuser.UserRefreshTokenRequest)
	if err := c.BodyParser(userRefreshTokenRequest); err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshTokenError),
			err.Error(),
		).Response()
	}

	var err = h.userUsecase.UserLogout(userRefreshTokenRequest)
	if err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(logoutUserError),
			err.Error(),
		).Response()
	}
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, nil).Response()
}

// RefreshToken refreshes user's authentication token.
// @Summary Refresh User Token
// @Description Refresh User Token
// @Tags User
// @Accept application/json
// @Produce json
// @Param request body entityuser.UserRefreshTokenRequest true "query params"
// @Success 200 {object} entityuser.UserLoginResponse
// @Failure 400 {object} gofiberentities.ErrorResponseType
// @Router /api/v1/refresh-token [post]
func (h *userHandler) RefreshToken(c *fiber.Ctx) error {
	var userRefreshTokenRequest = new(entityuser.UserRefreshTokenRequest)
	if err := c.BodyParser(userRefreshTokenRequest); err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshTokenError),
			err.Error(),
		).Response()
	}

	var userLoginReponse, err = h.userUsecase.UserRefreshToken(userRefreshTokenRequest)
	if err != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshTokenError),
			err.Error(),
		).Response()
	}
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, userLoginReponse).Response()
}

// GetUserProfile retrieves the user's profile.
// @Summary Get User Profile
// @Description Get User Profile
// @Tags User
// @Accept application/json
// @Produce json
// @Success 200 {object} entityuser.UserEntity
// @Failure 401 {object} gofiberentities.ErrorResponseType
// @Router /api/v1/profile [get]
func (h *userHandler) GetUserProfile(c *fiber.Ctx) error {
	var userEntity = c.Locals("user")
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, userEntity).Response()
}
