package repositories

import entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"

type ICRUDRepository interface {
	// Get
	GetOneById(id string) (*entities.DataResponse, error)

	// List
	GetList(pageNumber int, pageSize int, filter map[string]interface{}) (*entities.DataListResponse, error)

	// Create
	Create(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error)

	// Update
	Update(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error)

	// UpSert
	UpSert(filter map[string]interface{}, entityModel entities.ICRUDDataModel) (interface{}, error)

	// Delete
	Delete(entityModel *entities.ICRUDDataModel) error

	// Describe
	DescribeDataSource(schema string) map[string]interface{}
}
