package mongodbhandler

import (
	"context"

	"github.com/max38/golang-clean-code-architecture/src/config"
	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongodbHandler interface {
	GetDatabase(databaseName string) *mongo.Database
	GetCorrection(dataModel entities.ICRUDDataModel) *mongo.Collection
	GetCorrectionDB(database *mongo.Database, dataModel entities.ICRUDDataModel) *mongo.Collection
	GetContext() context.Context
	Close() error
}

type mongodbhandler struct {
	Client   *mongo.Client
	Database *mongo.Database
	Context  context.Context
}

func MongodbHandlerDefaultConfig() IMongodbHandler {
	return MongodbHandler(
		config.Config.String("GOAPP_MONGODB_DB_ATLAS_URI"),
		config.Config.String("GOAPP_MONGODB_DB_NAME"),
	)
}

func MongodbHandler(atlasUri string, databaseName string) IMongodbHandler {

	// var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	var ctx = context.Background()
	var client, err = mongo.Connect(ctx, options.Client().ApplyURI(atlasUri))

	if err != nil {
		panic(err)
	}
	return &mongodbhandler{
		Client:   client,
		Database: client.Database(databaseName),
		Context:  ctx,
	}

}

func (m *mongodbhandler) GetDatabase(databaseName string) *mongo.Database {
	return m.Client.Database(databaseName)
}

func (m *mongodbhandler) GetCorrection(dataModel entities.ICRUDDataModel) *mongo.Collection {
	return m.Database.Collection(dataModel.TableName())
}

func (m *mongodbhandler) GetCorrectionDB(database *mongo.Database, dataModel entities.ICRUDDataModel) *mongo.Collection {
	return database.Collection(dataModel.TableName())
}

func (m *mongodbhandler) GetContext() context.Context {
	return m.Context
}

func (m *mongodbhandler) Close() error {
	return m.Client.Disconnect(context.Background())
}
