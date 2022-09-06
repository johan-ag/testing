package users

import (
	"context"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/johan-ag/testing/internal/platform/database"
)

//---go:generate mockgen -destination=mocks/repository.go -package=mocks github.com/johan-ag/testing/internal/users Repository
type Repository interface {
	Save(ctx context.Context, name string, age uint, random string) (uint, error)
	Find(ctx context.Context, id uint) (User, error)
}

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func NewRepository(queries *database.Queries) *repository {
	return &repository{
		queries,
	}
}

type repository struct {
	queries *database.Queries
}

func (r *repository) Save(ctx context.Context, name string, age uint, random string) (uint, error) {
	result, err := r.queries.SaveUser(ctx, database.SaveUserParams{
		Name:   name,
		Age:    int32(age),
		Random: random,
	})
	if err != nil {
		return 0, ErrorSavingToDB
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, ErrorFindLastInsertedID
	}

	return uint(lastInsertID), nil
}

func (r *repository) Find(ctx context.Context, id uint) (User, error) {
	u, err := r.queries.FindUser(ctx, int32(id))
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:   uint(u.ID),
		Name: u.Name,
		Age:  uint(u.Age),
	}

	return user, nil
}
