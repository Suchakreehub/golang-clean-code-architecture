package gofibercrudhandler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
	gofiberentities "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/entities"
	crudusecase "github.com/max38/golang-clean-code-architecture/src/usecases/crud"
)

type ICRUDHandler interface {
	// Get
	Retrive(c *fiber.Ctx) error
	RetriveList(c *fiber.Ctx) error
	Describe(c *fiber.Ctx) error
}

type crudHandler struct {
	// crudUsecase *crudusecase.ICRUDUsecase
}

func CRUDHandler() ICRUDHandler {
	// var crudUsecase = crudusecase.CrudUsecase()
	return &crudHandler{
		// crudUsecase: &crudUsecase,
	}
}

func (h *crudHandler) Retrive(c *fiber.Ctx) error {
	var entityModelSlug = strings.Trim(c.Params("entity_model_slug"), " ")
	var entityId = c.Params("id")

	var modifierUserPermission = c.Locals("user_permission").(*entitymodels.UserPermissionModel)

	var crudUsecase = crudusecase.CrudUsecase()
	var responseData, errorRetrive = crudUsecase.Retrive(modifierUserPermission, entityModelSlug, entityId)

	if errorRetrive != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"error",
			errorRetrive.Error(),
		).Response()
	}
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, responseData).Response()
}

func (h *crudHandler) RetriveList(c *fiber.Ctx) error {
	var entityModelSlug = strings.Trim(c.Params("entity_model_slug"), " ")

	var modifierUserPermission = c.Locals("user_permission").(*entitymodels.UserPermissionModel)

	queryParams := make(map[string]interface{})
	for key, value := range c.Queries() {
		queryParams[key] = value
	}

	var page = 1
	if queryParams["_page"] != nil {
		if intValue, err := strconv.Atoi(queryParams["_page"].(string)); err == nil {
			page = intValue
			delete(queryParams, "_page")
		}
	}
	var pageSize = 10
	if queryParams["_pageSize"] != nil {
		if intValue, err := strconv.Atoi(queryParams["_pageSize"].(string)); err == nil {
			pageSize = intValue
			delete(queryParams, "_pageSize")
		}
	}

	var crudUsecase = crudusecase.CrudUsecase()
	var responseData, errorRetrive = crudUsecase.RetriveList(modifierUserPermission, entityModelSlug, page, pageSize, nil)

	if errorRetrive != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"error",
			errorRetrive.Error(),
		).Response()
	}
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, responseData).Response()
}

func (h *crudHandler) Describe(c *fiber.Ctx) error {
	var entityModelSlug = strings.Trim(c.Params("entity_model_slug"), " ")
	var schema = strings.Trim(c.Params("schema"), " ")
	var crudUsecase = crudusecase.CrudUsecase()
	var responseData, errorRetrive = crudUsecase.DescribeDataSource(entityModelSlug, schema)

	if errorRetrive != nil {
		return gofiberentities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"error",
			errorRetrive.Error(),
		).Response()
	}
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, responseData).Response()
}
