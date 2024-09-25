package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/handler"
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/SanyaWarvar/temple_api/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	dbNum, err := strconv.Atoi(os.Getenv("CACHE_DB"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}
	redisOptions := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PORT")),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       dbNum,
	}
	codeExp, err := time.ParseDuration(os.Getenv("CODE_EXP"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}
	cacheDb, err := repository.NewRedisDb(&redisOptions)
	if err != nil {
		logrus.Fatalf("Error while create connection to cache: %s", err.Error())
	}

	repos := repository.NewRepository(db, cacheDb, codeExp)
	tokenTTL, err := time.ParseDuration(os.Getenv("TOKENTTL"))
	if err != nil {
		logrus.Fatalf("Errof while parse tokenttl: %s", err.Error())
	}
	authSettings := service.NewAuthSettings(tokenTTL, os.Getenv("SALT"), os.Getenv("SIGNINGKEY"))

	codeLenght, err := strconv.Atoi(os.Getenv("CODE_LENGHT"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}

	smtpSettings := service.NewEmailSettings(os.Getenv("OWNER_EMAIL"), os.Getenv("OWNER_PASSWORD"), os.Getenv("SMTP_ADDR"), codeLenght)
	services := service.NewService(repos, *authSettings, *smtpSettings)

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
