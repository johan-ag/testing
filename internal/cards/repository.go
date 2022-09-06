package books

import (
	"context"
)

type Repository interface {
	Save(title string, author string) error
	Find(id string) (Book, error)
}

type Book struct {
}

type repository struct{}

func (r *repository) Save(ctx context.Context, title string, author string) error {
	return nil
}

func (r *repository) Find(id uint) {

}
