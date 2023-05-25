package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"selfuelAPI/internal/models"
)

type (
	MovieAPIProxy struct {
		apiKey  string
		token   string
		baseURL string
	}

	Movie struct {
		ID     int                   `json:"id"`
		Title  string                `json:"title"`
		Year   string                `json:"year"`
		Genres []*models.MovieGenres `json:"genres"`
		// TODO: fill the struct with rest of the response from the api call
	}
)

func NewMovieAPIProxy(apiKey, token, baseURL string) *MovieAPIProxy {
	return &MovieAPIProxy{
		apiKey:  apiKey,
		token:   token,
		baseURL: baseURL,
	}
}

func (mp *MovieAPIProxy) GetMovieData(ctx context.Context, movieID string) (*Movie, error) {
	parameters := url.Values{}
	parameters.Add("api_key", mp.apiKey)

	u, err := url.JoinPath(mp.baseURL, movieID, parameters.Encode())
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", mp.token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	movie := &Movie{}
	if err := json.Unmarshal(body, movie); err != nil {
		return nil, err
	}

	return movie, nil
}
