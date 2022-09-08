//go:generate mockgen -destination=./mocks.go -package=users github.com/johan-ag/testing/internal/users . Service,Repository
//go:generate mockgen -destination=../../internal/platform/kvs/mock.go -package=kvs github.com/mercadolibre/fury_go-toolkit-kvs/pkg/kvs QueryableClient
package users

import (
	"context"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository,
	}
}

type Service interface {
	Save(ctx context.Context, name string, age uint) (User, error)
	Find(ctx context.Context, id uint) (User, error)
	FindByParams(ctx context.Context, name string, age int32) ([]User, error)
}

func (s *service) Save(ctx context.Context, name string, age uint) (User, error) {
	id, err := s.repository.Save(ctx, name, age)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:   id,
		Name: name,
		Age:  age,
	}

	return user, nil
}

func (s *service) Find(ctx context.Context, id uint) (User, error) {
	return s.repository.Find(ctx, id)
}

func (s *service) FindByParams(ctx context.Context, name string, age int32) ([]User, error) {
	return s.repository.FindByNameAndAge(ctx, name, age)
}
