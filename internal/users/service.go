package users

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/mercadolibre/fury_go-toolkit-kvs/pkg/kvs"
)

type Service interface {
	Save(ctx context.Context, name string, age uint) (uint, error)
	Find(ctx context.Context, id uint) (User, error)
}

//go:generate mockgen -destination=./mocks.go -package=users github.com/johan-ag/testing/internal/users Repository,Service
//go:generate mockgen -destination=../../internal/platform/kvs/mock.go -package=kvs github.com/mercadolibre/fury_go-toolkit-kvs/pkg/kvs QueryableClient
type service struct {
	repository Repository
	qkvs       kvs.QueryableClient
}

func NewService(repository Repository, qkvs kvs.QueryableClient) *service {
	return &service{
		repository,
		qkvs,
	}
}

// Save method save
func (s *service) Save(ctx context.Context, name string, age uint) (uint, error) {
	random, err := generateRandom()
	if err != nil {
		return 0, err
	}

	id, err := s.repository.Save(ctx, name, age, random)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) Find(ctx context.Context, id uint) (User, error) {
	user, err := s.repository.Find(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// generateRandom generate length six random string using go-nanoid library,
func generateRandom() (string, error) {
	activationAlphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" //TODO
	activationLength := 6

	random, err := gonanoid.Generate(activationAlphabet, activationLength)
	if err != nil {
		return "", err
	}

	return random, nil
}
