package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	usersHandler "github.com/johan-ag/testing/cmd/api/users"
	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-core/pkg/log"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
)

func main() {
	if err := run(); err != nil {
		log.Error(context.Background(), "cannot run application", log.Err(err))
	}
}
func run() error {
	app, err := fury.NewWebApplication()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", "root:root@/testdb")
	if err != nil {
		return err
	}

	usersRepository := users.NewRepository(db)
	usersService := users.NewService(usersRepository)

	usersHandler := usersHandler.NewHandler(usersService)

	app.Post("/api/users", usersHandler.Save)
	app.Get("/api/users/{id}", usersHandler.Find)
	app.Get("/api/users", usersHandler.FindByNameAndAge)

	return app.Run()
}
