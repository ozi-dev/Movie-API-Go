package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri       = "mongodb+srv://<your info>@movie-details.attlleo.mongodb.net/?retryWrites=true&w=majority"
	apiKey    = "<apikey>"
	apiURL    = "https://api.themoviedb.org/3/movie"
	apiBearer = "Bearer 1"
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
	http.HandleFunc("/movies", handleMoviesRequest)

	// Sunucuyu başlat
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func handleMoviesRequest(w http.ResponseWriter, r *http.Request) {
	movieID := r.URL.Query().Get("id")

	if movieID == "" {
		http.Error(w, "Missing movie ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/%s?api_key=%s", apiURL, movieID, apiKey)
	method := "GET"

	api := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", apiBearer)

	res, err := api.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiResponse := string(body)

	var movie Movie
	err = json.Unmarshal([]byte(apiResponse), &movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var foundMovies []Movie
	for cursor.Next(ctx) {
		var movie Movie
		err := cursor.Decode(&movie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		foundMovies = append(foundMovies, movie)
	}

	if foundMovies == nil {
		_, err := collection.InsertMany(ctx, movies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Inserted movie details:\nID: %s\nTitle: %s\nGenres: %v\nYear: %s", movieID, title, genres, year)
		w.Write([]byte(response))

	} else {
		response := fmt.Sprintf("Movie found in database:\nID: %s\nTitle: %s\nGenres: %v\nYear: %s", movieID, title, genres, year)
		w.Write([]byte(response))
	}
}
