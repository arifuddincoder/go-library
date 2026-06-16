package main

import (
	"go-library/internal/config"
	"go-library/internal/domain/user"
	"go-library/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	user.SeedAdmin(user.NewRepository(db), cfg)
	server.Start(db, cfg)
}
