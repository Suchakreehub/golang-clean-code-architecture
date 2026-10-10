package entityuser

import "time"

type UserEntity struct {
	Id        uint      `db:"id" gorm:"primaryKey" json:"id"`
	Email     string    `db:"email" gorm:"uniqueIndex" json:"email"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserTokenEntity struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type UserTokenClaimsEntity struct {
	Id uint `json:"id"`
}
