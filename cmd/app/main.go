package main

import (
	"os"

	"github.com/SanyaWarvar/temple_api/pkg/handler"
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/SanyaWarvar/temple_api/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error while load dotenv: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	handlers := handler.NewHandler(services)
	srv := new(models.Server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error while running server: %s", err.Error())
	}

}
