package entities

import (
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"
)

type ICRUDDataModel interface {
	TableName() string
	Datasource() sharedcrud.DatasourceType
}
