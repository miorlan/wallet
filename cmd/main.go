package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
	"wallet/config"
	"wallet/internal/api"
	"wallet/internal/repository"
	"wallet/internal/service"
)

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	repo := repository.NewDBWallet(db)
	svc := service.NewWalletService(repo)
	handler := api.NewWalletHandler(svc)

	r := handler.SetupRoutes()

	log.Printf("Server started on %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, r))
}
