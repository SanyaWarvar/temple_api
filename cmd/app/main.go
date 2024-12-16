package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/handler"
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/SanyaWarvar/temple_api/pkg/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	err := os.MkdirAll("user_data/profile_pictures", 0750)
	err = os.MkdirAll("user_data/tik_toks", 0750)

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

	//err = generateStatics(db)
	if err != nil {
		logrus.Fatalf("Error while create statics: %s", err.Error())
	}

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
	codeLenght, err := strconv.Atoi(os.Getenv("CODE_LENGHT"))
	if err != nil {
		logrus.Fatalf("Error while create connection to cache: %s", err.Error())
	}
	emailCfg := repository.NewEmailCfg(os.Getenv("OWNER_EMAIL"), os.Getenv("OWNER_PASSWORD"), os.Getenv("SMTP_ADDR"), codeLenght)

	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESSTOKENTTL"))
	if err != nil {
		logrus.Fatalf("Errof while parse accessTokenTTL: %s", err.Error())
	}
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESHTOKENTTL"))
	if err != nil {
		logrus.Fatalf("Errof while parse refreshTokenTTL: %s", err.Error())
	}
	jwtCfg := repository.NewJwtManagerCfg(accessTokenTTL, refreshTokenTTL, os.Getenv("SIGNINGKEY"), jwt.SigningMethodHS256)

	repos := repository.NewRepository(db, cacheDb, codeExp, emailCfg, jwtCfg)

	services := service.NewService(repos)

	handlers := handler.NewHandler(services)
	srv := new(models.Server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	releaseMode, err := strconv.ParseBool(os.Getenv("RELEASEMODE"))
	if err != nil {
		logrus.Fatalf("Error while parse Realease mode from .env: %s", err.Error())
	}
	if err := srv.Run(port, handlers.InitRoutes(releaseMode)); err != nil {
		logrus.Fatalf("Error while running server: %s", err.Error())
	}

}

type StaticFile struct {
	Filename     string `db:"filename"`
	FileAsString string `db:"file"`
	File         []byte
}

func generateStatics(db *sqlx.DB) error {
	var files []StaticFile

	query := `
		SELECT (select username from users where id = ui.user_id) as filename, profile_picture as file FROM users_info ui
	`
	err := db.Select(&files, query)
	if err != nil {
		return err
	}
	for ind, item := range files {

		files[ind].File, err = base64.RawStdEncoding.DecodeString(item.FileAsString)
		if err != nil {
			continue
		}
		os.WriteFile("user_data/profile_pictures/"+files[ind].Filename, files[ind].File, 0755)
	}
	files = []StaticFile{}
	query = `
		SELECT id as filename, body as file FROM tiktoks
	`
	err = db.Select(&files, query)
	if err != nil {
		return err
	}
	for ind, item := range files {
		files[ind].File, err = base64.RawStdEncoding.DecodeString(item.FileAsString)
		if err != nil {
			continue
		}
		os.WriteFile("user_data/tik_toks/"+files[ind].Filename, files[ind].File, 0755)
	}
	return nil
}
