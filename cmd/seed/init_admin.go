package main

import (
	"flag"
	"log"

	"github.com/ximofam/user-service/internal/config"
	"github.com/ximofam/user-service/internal/model"
	"github.com/ximofam/user-service/internal/utils"
	"github.com/ximofam/user-service/pkg/datasource"
)

func main() {
	// 1. Định nghĩa các tham số (Flags)
	email := flag.String("email", "", "Email của Admin")
	password := flag.String("password", "", "Mật khẩu của Admin")
	username := flag.String("username", "admin", "Username (mặc định là admin)")

	flag.Parse()

	if *email == "" || *password == "" || *username == "" {
		log.Fatal("Failed to init admin: admin must be have enough 3 things: username, email, password")
	}

	cfg := config.Load()

	db, err := datasource.NewMySQLConnection(
		cfg.MySQL.DBUser,
		cfg.MySQL.DBPassword,
		cfg.MySQL.DBHost,
		cfg.MySQL.DBPort,
		cfg.MySQL.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to init admin: %v", err)
	}

	var count int
	db.QueryRow("SELECT 1 FROM users WHERE email = ? OR username = ? LIMIT 1",
		email,
		username).
		Scan(&count)

	if count != 0 {
		log.Fatalf("Failed to init admin: Admin already exists with email=%s or username=%s", *email, *username)
	}

	hash := utils.HashPassword(*password)

	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES(?,?,?,?)",
		email, username, hash, model.RoleAdmin,
	)

	if err != nil {
		log.Fatalf("Failed to init admin: %v", err)
	}

	log.Printf("Created admin with email=%s, username=%s, password=%s\n", *email, *username, *password)
}
