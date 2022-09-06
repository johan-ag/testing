package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	usersHandler "github.com/johan-ag/testing/cmd/api/users"
	"github.com/johan-ag/testing/internal/platform/database"
	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-core/pkg/log"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/mercadolibre/fury_go-toolkit-kvs/pkg/kvs"
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

	queries := database.New(db)

	qkvs, err := kvs.NewQueryableClient("container")
	if err != nil {
		return err
	}
	usersRepository := users.NewRepository(queries)
	usersService := users.NewService(usersRepository, qkvs)

	_usersHandler := usersHandler.NewHandler(usersService)

	app.Post("/api/users", _usersHandler.Save)
	app.Get("/api/users/{id}", _usersHandler.Find)

	return app.Run()
}
