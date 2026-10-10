package crudusecase

import (
	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	entitymodels "github.com/max38/golang-clean-code-architecture/src/domain/models"
	crudrepository "github.com/max38/golang-clean-code-architecture/src/interface/repositories/crud"
)

type ICRUDUsecase interface {
	// Get
	Retrive(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, id string) (interface{}, error)

	// List
	RetriveList(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, pageNumber int, pageSize int, filter map[string]interface{}) (*entities.DataListResponse, error)

	// Create
	Create(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) (interface{}, error)

	// Update
	Update(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) (interface{}, error)

	// Delete
	Delete(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) error

	// Describe
	DescribeDataSource(entityModelSlug string, schema string) (map[string]interface{}, error)
}

type crudUsecase struct {
}

func CrudUsecase() ICRUDUsecase {
	return &crudUsecase{}
}

func (c *crudUsecase) Retrive(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, id string) (interface{}, error) {
	var crudRepository, errorInit = crudrepository.GetCRUDRepository(entityModelSlug)
	if errorInit != nil {
		return nil, errorInit
	}
	return crudRepository.GetOneById(id)
}

func (c *crudUsecase) RetriveList(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, pageNumber int, pageSize int, filter map[string]interface{}) (*entities.DataListResponse, error) {
	var crudRepository, errorInit = crudrepository.GetCRUDRepository(entityModelSlug)

	if errorInit != nil {
		return nil, errorInit
	}
	return crudRepository.GetList(pageNumber, pageSize, filter)
}

func (c *crudUsecase) Create(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) (interface{}, error) {
	return nil, nil
	// var crudRepository, errorInit = crudrepository.GetCRUDRepository(entityModelSlug)
	// if errorInit != nil {
	// 	return nil, errorInit
	// }
	// return crudRepository.Create(entityModel)
}

func (c *crudUsecase) Update(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) (interface{}, error) {
	return nil, nil
}

func (c *crudUsecase) Delete(userPermissionModel *entitymodels.UserPermissionModel, entityModelSlug string, entityModel interface{}) error {
	return nil
}

func (c *crudUsecase) DescribeDataSource(entityModelSlug string, schema string) (map[string]interface{}, error) {
	var crudRepository, errorInit = crudrepository.GetCRUDRepository(entityModelSlug)

	if errorInit != nil {
		return nil, errorInit
	}
	return crudRepository.DescribeDataSource(schema), nil
}
