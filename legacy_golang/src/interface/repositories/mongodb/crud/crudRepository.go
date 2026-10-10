package mongodbcrudrepository

import (
	"math"
	"reflect"
	"time"

	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	"github.com/max38/golang-clean-code-architecture/src/domain/repositories"
	mongodbhandler "github.com/max38/golang-clean-code-architecture/src/infrastructure/database/mongodb"
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type crudRepository struct {
	dbHandler    mongodbhandler.IMongodbHandler
	entityModel  *entities.ICRUDDataModel
	dbCorrection *mongo.Collection
}

func CrudRepository(entityModel *entities.ICRUDDataModel) repositories.ICRUDRepository {
	var dbHandler = mongodbhandler.MongodbHandlerDefaultConfig()

	var dbCorrection = dbHandler.GetCorrection(*entityModel)
	return &crudRepository{
		dbHandler:    dbHandler,
		entityModel:  entityModel,
		dbCorrection: dbCorrection,
	}
}

func (c *crudRepository) convertToBson(source map[string]interface{}) (bson.M, error) {
	var bsonMap = bson.M{}
	for key, value := range source {
		bsonMap[key] = value
	}
	return bsonMap, nil
}

// Get
func (c *crudRepository) GetOneById(id string) (*entities.DataResponse, error) {
	// var entityModel = reflect.New(reflect.ValueOf(*c.entityModel).Type()).Interface()
	var entityModel = reflect.New(reflect.TypeOf(*c.entityModel)).Interface()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var filter = bson.M{"_id": objectID}
	var dbResult = c.dbCorrection.FindOne(c.dbHandler.GetContext(), filter).Decode(entityModel)
	if dbResult != nil {
		return nil, dbResult
	}

	var dataResponse = entities.DataResponse{
		Data:       entityModel,
		DataSchema: c.DescribeDataSource("detail"),
	}

	return &dataResponse, nil
}

// List
func (c *crudRepository) GetList(pageNumber int, pageSize int, filter map[string]interface{}) (*entities.DataListResponse, error) {
	if pageSize == 0 {
		pageSize = 10
	} else if pageSize > 1000 {
		pageSize = 1000
	}

	if pageNumber < 1 {
		pageNumber = 1
	}

	var entityType = reflect.ValueOf(*c.entityModel).Type()
	var entityModels = reflect.MakeSlice(reflect.SliceOf(entityType), 0, 0).Interface()

	var filterData, errorConvert = c.convertToBson(filter)
	if errorConvert != nil {
		return nil, errorConvert
	}
	var offsetSkip = (pageNumber - 1) * pageSize

	var options = options.Find()
	options.SetSort(bson.M{"_id": 1})
	options.SetSkip(int64(offsetSkip))
	options.SetLimit(int64(pageSize))

	cursor, err := c.dbCorrection.Find(c.dbHandler.GetContext(), filterData, options)

	if err != nil {
		return nil, err
	}
	if err = cursor.All(c.dbHandler.GetContext(), &entityModels); err != nil {
		return nil, err
	}

	var totalCount, errCount = c.dbCorrection.CountDocuments(c.dbHandler.GetContext(), filterData)
	if errCount != nil {
		return nil, errCount
	}

	var dataResponse = entities.DataListResponse{
		Data:       entityModels,
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPage:  int(math.Ceil(float64(totalCount) / float64(pageSize))),
		TotalData:  totalCount,
		DataSchema: nil,
	}

	return &dataResponse, nil
}

// Create
func (c *crudRepository) Create(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error) {
	return nil, nil
}

// Update
func (c *crudRepository) Update(entityModel *entities.ICRUDDataModel) (*entities.DataResponse, error) {
	return nil, nil
}

func (c *crudRepository) UpSert(filter map[string]interface{}, entityModelUpdate entities.ICRUDDataModel) (interface{}, error) {
	var opts = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var update = bson.D{
		{"$set", entityModelUpdate},
		{"$currentDate", bson.D{{"updated_at", true}}},
		{"$setOnInsert", bson.D{{"created_at", primitive.NewDateTimeFromTime(time.Now())}}},
	}

	var filterData, errorConvert = c.convertToBson(filter)
	if errorConvert != nil {
		return nil, errorConvert
	}

	var entityModelUpdated = reflect.New(reflect.TypeOf(entityModelUpdate)).Interface()
	var err = c.dbCorrection.FindOneAndUpdate(c.dbHandler.GetContext(), filterData, update, opts).Decode(&entityModelUpdated)
	if err != nil {
		return nil, err
	}

	return entityModelUpdated, nil
}

// Delete
func (c *crudRepository) Delete(entityModel *entities.ICRUDDataModel) error {
	return nil
}

// Delete
func (c *crudRepository) DescribeDataSource(schema string) map[string]interface{} {
	var entityModel = reflect.New(reflect.ValueOf(*c.entityModel).Type().Elem()).Interface()
	var schemaDetail = sharedcrud.GenerateAvroSchema(entityModel)
	return schemaDetail
}
