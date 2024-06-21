package test

import (
	"dealls-dating/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var log *logrus.Logger

var validate *validator.Validate

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	log = config.NewLog()
	app = config.NewFiber()
	db = config.NewDB(log)
	validate = config.NewValidate()

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Log:      log,
		Validate: validate,
	})
}
