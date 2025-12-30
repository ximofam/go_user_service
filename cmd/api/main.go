package main

import (
	"log"

	"github.com/ximofam/user-service/internal/app"
	"github.com/ximofam/user-service/internal/config"
	"github.com/ximofam/user-service/pkg/datasource"
)

func main() {
	cfg := config.Load()

	db, err := datasource.NewMySQLConnection(
		cfg.MySQL.DBUser,
		cfg.MySQL.DBPassword,
		cfg.MySQL.DBHost,
		cfg.MySQL.DBPort,
		cfg.MySQL.DBName,
	)

	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := datasource.NewRedisClient(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	if err != nil {
		log.Fatal(err)
	}

	app, err := app.NewApp(cfg, db, redisClient)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
