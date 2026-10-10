package entityuser

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserRegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *UserRegisterRequest) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, r.Email)
	if err != nil {
		return false
	}
	return match
}

func (r *UserRegisterRequest) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}
	r.Password = string(hashedPassword)
	return nil
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
