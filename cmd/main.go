package main

import (
	"go-library/internal/config"
	"go-library/internal/domain/book"
	"go-library/internal/domain/category"
	"go-library/internal/domain/user"
	"go-library/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	db.AutoMigrate(&user.User{}, &category.Category{}, &book.Book{})
	user.SeedAdmin(user.NewRepository(db), cfg)
	server.Start(db, cfg)
}
