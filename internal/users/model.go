//go:generate mockgen -destination=mocks.go -package=users github.com/johan-ag/testing/internal/users . Repository
package users

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}
