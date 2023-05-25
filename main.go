package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri       string
	apiKey    string
	apiURL    string
	apiBearer string
)

type Movie struct {
	ID     int `json:"id"`
	Title  string
	Genres []struct {
		Name string `json:"name"`
	}
	Year string `json:"release_date"`
}

func main() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// .env dosyasından değerleri al
	uri = os.Getenv("URI")
	apiKey = os.Getenv("API_KEY")
	apiURL = os.Getenv("API_URL")
	apiBearer = os.Getenv("API_BEARER")

	e := echo.New()

	// Middleware'leri ekle
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Endpoint'leri tanımla
	e.GET("/movies", handleMoviesRequest)

	// Sunucuyu başlat
	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}

func handleMoviesRequest(c echo.Context) error {
	movieID := c.QueryParam("id")

	if movieID == "" {
		return c.String(http.StatusBadRequest, "Missing movie ID")
	}

	url := fmt.Sprintf("%s/%s?api_key=%s", apiURL, movieID, apiKey)
	method := "GET"

	api := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Add("Authorization", apiBearer)

	res, err := api.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	apiResponse := string(body)

	var movie Movie
	err = json.Unmarshal([]byte(apiResponse), &movie)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testDatabase").Collection("testCollection")

	ID := movie.ID
	genres := movie.Genres
	title := movie.Title
	year := movie.Year[:4]

	var movies []interface{}
	movies = append(movies, Movie{ID: ID, Title: title, Genres: genres, Year: year})

	findFilter := bson.D{
		bson.E{Key: "id", Value: ID},
	}

	cursor, err := collection.Find(ctx, findFilter)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(ctx)

	var foundMovies []Movie
	for cursor.Next(ctx) {
		var movie Movie
		err := cursor.Decode(&movie)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		foundMovies = append(foundMovies, movie)
	}

	if foundMovies == nil {
		_, err := collection.InsertMany(ctx, movies)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		response := fmt.Sprintf("Inserted movie details:\nID: %s\nTitle: %s\nGenres: %v\nYear: %s", movieID, title, genres, year)
		return c.String(http.StatusOK, response)

	} else {
		response := fmt.Sprintf("Movie found in database:\nID: %s\nTitle: %s\nGenres: %v\nYear: %s", movieID, title, genres, year)
		return c.String(http.StatusOK, response)
	}
}
