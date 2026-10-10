package crudrepository

import (
	"fmt"

	config "github.com/max38/golang-clean-code-architecture/src/config"
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"

	entities "github.com/max38/golang-clean-code-architecture/src/domain/entities/crud"
	repositories "github.com/max38/golang-clean-code-architecture/src/domain/repositories"

	mongodbcrudrepository "github.com/max38/golang-clean-code-architecture/src/interface/repositories/mongodb/crud"
	postgrescrudrepository "github.com/max38/golang-clean-code-architecture/src/interface/repositories/postgres/crud"
)

func GetCRUDRepository(entityModelSlug string) (repositories.ICRUDRepository, error) {

	var entityModelMap = make(map[string]entities.ICRUDDataModel)

	for _, entityModel := range config.EntitiyModels {
		var tableName = entityModel.TableName()
		entityModelMap[sharedcrud.ConvertNameToCRUDSlug(tableName)] = entityModel
	}

	if entityModelUse, exists := entityModelMap[entityModelSlug]; exists {
		// Key exists in the map
		if entityModelUse.Datasource() == sharedcrud.DatasourcePostgresql {
			var crudRepository = postgrescrudrepository.CrudRepository(&entityModelUse)
			return crudRepository, nil
		} else if entityModelUse.Datasource() == sharedcrud.DatasourceMongodb {
			var crudRepository = mongodbcrudrepository.CrudRepository(&entityModelUse)
			return crudRepository, nil
		} else {
			return nil, fmt.Errorf("Datasource not supported")
		}
	}

	return nil, fmt.Errorf("Entity model not found")
}
