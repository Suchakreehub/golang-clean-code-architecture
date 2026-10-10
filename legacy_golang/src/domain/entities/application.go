package entities

import (
	userusecase "github.com/max38/golang-clean-code-architecture/src/usecases/user"
)

type ApplicationEntity struct {
	UserUsecase userusecase.IUserUsecase
}
