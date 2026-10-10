package postgrescrudrepository

import (
	"math"
	"reflect"

	"gorm.io/gorm"

	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	"github.com/max38/golang-clean-code-architecture/src/domain/repositories"

	postgreshandler "github.com/max38/golang-clean-code-architecture/src/infrastructure/database/postgres"
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"
)

type crudRepository struct {
	dbHandler   postgreshandler.IPostgresHandler
	entityModel *entities.ICRUDDataModel
}

func CrudRepository(entityModel *entities.ICRUDDataModel) repositories.ICRUDRepository {
	var dbHandler = postgreshandler.PostgresHandler()
	return &crudRepository{
		dbHandler:   dbHandler,
		entityModel: entityModel,
	}
}

func (c *crudRepository) getTableName(entityModel interface{}) string {
	newEntityModel := reflect.New(reflect.TypeOf(entityModel).Elem().Elem()).Interface()
	if obj, ok := newEntityModel.(interface{ TableName() string }); ok {
		return obj.TableName()
	}
	return ""
}

func (c *crudRepository) getQueryChain(entityModel interface{}) *gorm.DB {
	var queryChain = c.dbHandler.GetConnector().Model(entityModel)
	newEntityModel := reflect.New(reflect.TypeOf(entityModel).Elem().Elem()).Interface()
	if obj, ok := newEntityModel.(interface{ JoinAssociations() []string }); ok {
		var joinEntities = obj.JoinAssociations()

		for _, joinEntity := range joinEntities {
			queryChain = queryChain.Joins(string(joinEntity))
		}
	}

	return queryChain
}

// Get
func (c *crudRepository) GetOneById(id string) (*entities.DataResponse, error) {
	var entityModel = reflect.New(reflect.ValueOf(*c.entityModel).Type()).Interface()

	var dbConnector = c.getQueryChain(entityModel)
	var tableName = c.getTableName(entityModel)

	var dbResult = dbConnector.Where(tableName+".id = ?", id).First(entityModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	var dataResponse = entities.DataResponse{
		Data:       entityModel,
		DataSchema: c.DescribeDataSource("detail"),
	}
	return &dataResponse, nil
}

// List
func (c *crudRepository) GetList(pageNumber int, pageSize int, filter map[string]interface{}) (*entities.DataListResponse, error) {

	var entityType = reflect.ValueOf(*c.entityModel).Type()
	var entityModel = reflect.New(entityType).Interface()
	var dbConnector = c.getQueryChain(entityModel)

	var entityModels = reflect.MakeSlice(reflect.SliceOf(entityType), 0, 0).Interface()

	if pageSize == 0 {
		pageSize = 10
	} else if pageSize > 1000 {
		pageSize = 1000
	}

	if pageNumber < 1 {
		pageNumber = 1
	}

	var dbResult = dbConnector.Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&entityModels)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	var dataResponse = entities.DataListResponse{
		Data:       entityModels,
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPage:  int(math.Ceil(float64(dbResult.RowsAffected) / float64(pageSize))),
		TotalData:  dbResult.RowsAffected,
		DataSchema: c.DescribeDataSource("list"),
	}

	return &dataResponse, nil
}

// Create
func (c *crudRepository) Create(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error) {
	var dbConnector = c.dbHandler.GetConnector()
	var dbResult = dbConnector.Create(&entityModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	var dataResponse = entities.DataResponse{
		Data:       entityModel,
		DataSchema: nil,
	}
	return &dataResponse, nil
}

// Update
func (c *crudRepository) Update(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error) {
	var dbConnector = c.dbHandler.GetConnector()
	var dbResult = dbConnector.Save(&entityModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	var dataResponse = entities.DataResponse{
		Data:       entityModel,
		DataSchema: nil,
	}
	return &dataResponse, nil
}

// UpSert
func (c *crudRepository) UpSert(filter map[string]interface{}, entityModel entities.ICRUDDataModel) (interface{}, error) {
	var dbConnector = c.dbHandler.GetConnector()
	var dbResult = dbConnector.Where(filter).Assign(entityModel).FirstOrCreate(&entityModel)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &entityModel, nil
}

// Delete
func (c *crudRepository) Delete(entityModel *entities.ICRUDDataModel) error {
	var dbConnector = c.dbHandler.GetConnector()
	var dbResult = dbConnector.Delete(&entityModel)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

// Describe
func (c *crudRepository) DescribeDataSource(schema string) map[string]interface{} {
	// not finished yet

	var entityModel = reflect.New(reflect.ValueOf(*c.entityModel).Type().Elem()).Interface()
	var schemaDetail = sharedcrud.GenerateAvroSchema(entityModel)

	// var x = map[string]interface{}{
	// 	"r": schema,
	// }

	return schemaDetail
}

// Avro
