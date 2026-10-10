package main

import (
	"fmt"
	"log"

	"github.com/max38/golang-clean-code-architecture/src/config"
	postgreshandler "github.com/max38/golang-clean-code-architecture/src/infrastructure/database/postgres"
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"
)

func main() {
	config.Load(".env.dev")

	var dbHandler = postgreshandler.PostgresHandler()
	var dbConnector = dbHandler.GetConnector()

	log.Println("Connect success")
	log.Println("Start migration")

	// Migrate the schema
	for _, entityModel := range config.EntitiyModels {
		if entityModel.Datasource() == sharedcrud.DatasourcePostgresql {
			dbConnector.AutoMigrate(entityModel)
			fmt.Println("Migrate table: " + entityModel.TableName())
		}
	}

	dbHandler.Close()

	log.Println("Migration success !!!")
}
