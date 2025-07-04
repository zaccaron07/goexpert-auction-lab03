package user_entity

import (
	"context"

	"github.com/zaccaron07/goexpert-auction-lab03/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
