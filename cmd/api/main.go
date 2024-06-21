package main

import (
	"dealls-dating/internal/config"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log := config.NewLog()

	app := config.NewFiber()
	db := config.NewDB(log)
	validate := config.NewValidate()

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Log:      log,
		Validate: validate,
	})

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT"))); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
