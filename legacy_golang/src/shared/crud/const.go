package sharedcrud

type DatasourceType string

const (
	DatasourcePostgresql DatasourceType = "postgresql"
	DatasourceMongodb    DatasourceType = "mongodb"
)
