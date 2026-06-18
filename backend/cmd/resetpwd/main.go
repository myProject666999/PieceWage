package main

import (
	"fmt"

	"piece-wage/internal/config"
	"piece-wage/pkg/db"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		panic(err)
	}
	if err := db.Init(&cfg.MySQL); err != nil {
		panic(err)
	}

	password := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashStr := string(hash)
	fmt.Println("Generated hash:", hashStr)

	result := db.DB.Exec("UPDATE sys_user SET password = ?", hashStr)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Printf("Updated %d rows\n", result.RowsAffected)
}
