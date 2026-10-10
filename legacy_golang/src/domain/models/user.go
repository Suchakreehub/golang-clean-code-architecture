package entitymodels

import (
	"time"

	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
	sharedcrud "github.com/max38/golang-clean-code-architecture/src/shared/crud"
)

type UserModel struct {
	entityuser.UserEntity
	Password string `db:"password" json:"-"`
}

func (u *UserModel) TableName() string {
	return "users"
}

func (u *UserModel) ToEntity() *entityuser.UserEntity {
	return &u.UserEntity
}

func (u *UserModel) Datasource() sharedcrud.DatasourceType {
	return sharedcrud.DatasourcePostgresql
}

// ---------

type UserTokenModel struct {
	ID        uint      `gorm:"primary_key"`
	UserId    uint      `gorm:"index"`
	User      UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Access    string    `json:"access"`
	Refresh   string    `json:"refresh"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (u *UserTokenModel) TableName() string {
	return "user_tokens"
}

func (u *UserTokenModel) Datasource() sharedcrud.DatasourceType {
	return sharedcrud.DatasourcePostgresql
}

// ---------

type UserPermissionModel struct {
	ID         uint      `gorm:"primary_key"`
	UserId     uint      `gorm:"index"`
	User       UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Permission string    `json:"permission"` // You can edit this field
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (u *UserPermissionModel) TableName() string {
	return "user_permissions"
}

func (u *UserPermissionModel) Datasource() sharedcrud.DatasourceType {
	return sharedcrud.DatasourcePostgresql
}
