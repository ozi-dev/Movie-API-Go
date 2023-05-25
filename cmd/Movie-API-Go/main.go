package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"selfuelAPI/internal/api"
	"selfuelAPI/internal/repositories"
	"selfuelAPI/internal/server"
	"selfuelAPI/internal/service"
)

var (
	uri                       = os.Getenv("URI")
	apiKey                    = os.Getenv("API_KEY")
	apiURL                    = os.Getenv("API_URL")
	apiBearer                 = os.Getenv("DB_NAME")
	dbName                    = os.Getenv("API_BEARER")
	movieRepositoryCollection = os.Getenv("MOVIE_REPOSITORY_COLLECTION")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	defer client.Disconnect(ctx)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	movieAPIProxy := api.NewMovieAPIProxy(apiURL, apiKey, apiBearer)
	movieRepo := repositories.NewMovieRepository(client, dbName, movieRepositoryCollection)
	movieService := service.New(movieRepo, movieAPIProxy)
	movieServer := server.New(movieService)

	e.GET("/movies", movieServer.GetMovie)

	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}
