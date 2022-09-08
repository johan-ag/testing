package users

import (
	"context"
	"database/sql"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/johan-ag/testing/internal/platform/database"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks github.com/johan-ag/testing/internal/users Repository
type repository struct {
	queries  *database.Queries
	database *sql.DB
}

func NewRepository(queries *database.Queries, database *sql.DB) *repository {
	return &repository{
		queries:  queries,
		database: database,
	}
}

type Repository interface {
	Save(ctx context.Context, name string, age uint) (uint, error)
	Find(ctx context.Context, id uint) (User, error)
	FindByNameAndAge(ctx context.Context, name string, age int32) ([]User, error)
}

func (r *repository) Save(ctx context.Context, name string, age uint) (uint, error) {
	result, err := r.queries.SaveUser(ctx, database.SaveUserParams{
		Name: name,
		Age:  int32(age),
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

func (r *repository) FindByNameAndAge(ctx context.Context, name string, age int32) ([]User, error) {
	params := database.FindUserByParamsParams{
		Name: name,
		Age:  age,
	}

	dbUsers, err := r.queries.FindUserByParams(ctx, params)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, u := range dbUsers {
		user := User{
			ID:   uint(u.ID),
			Name: u.Name,
			Age:  uint(u.Age),
		}
		users = append(users, user)
	}

	return users, nil
}
