package postgreshandler

import (
	"fmt"

	"github.com/max38/golang-clean-code-architecture/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IPostgresHandler interface {
	GetConnector() *gorm.DB
	Close() error
}

type postgreshandler struct {
	DB *gorm.DB
}

func PostgresHandler() IPostgresHandler {
	var dsn string = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Config.String("GOAPP_POSTGRES_DB_HOST"),
		config.Config.String("GOAPP_POSTGRES_DB_USER"),
		config.Config.String("GOAPP_POSTGRES_DB_PASSWORD"),
		config.Config.String("GOAPP_POSTGRES_DB_NAME"),
		config.Config.Int("GOAPP_POSTGRES_DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return &postgreshandler{
		DB: db,
	}
}

func (p *postgreshandler) GetConnector() *gorm.DB {
	return p.DB
}

func (p *postgreshandler) Close() error {
	db, err := p.DB.DB()

	if err != nil {
		return err
	}

	return db.Close()
}
