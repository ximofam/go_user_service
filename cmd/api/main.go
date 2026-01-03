package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ximofam/user-service/internal/app"
	"github.com/ximofam/user-service/internal/config"
	"github.com/ximofam/user-service/internal/worker"
	"github.com/ximofam/user-service/pkg/datasource"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

	// Create sql.DB (using mysql)
	dbMySQL, err := datasource.NewMySQLConnection(
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

	worker.Init(3, 3)
	defer worker.Shutdown()

	app, err := app.NewApp(cfg, dbMySQL, redisClient)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
