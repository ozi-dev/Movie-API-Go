package service

import (
	"context"
	"errors"

	"selfuelAPI/internal/api"
	"selfuelAPI/internal/models"
	"selfuelAPI/internal/repositories"
)

type Service struct {
	movieRepository *repositories.MovieRepository
	movieAPIProxy   *api.MovieAPIProxy
}

func New(movieRepository *repositories.MovieRepository, movieAPIProxy *api.MovieAPIProxy) *Service {
	return &Service{
		movieRepository: movieRepository,
		movieAPIProxy:   movieAPIProxy,
	}
}

type (
	GetMovieRequest struct {
		ID string `query:"id"`
	}

	GetMovieResponse struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Genre string `json:"genre"`
		Year  string `json:"release_date"`
	}
)

func (req *GetMovieRequest) Validate() error {
	if len(req.ID) < 1 {
		return errors.New("id is required")
	}

	return nil
}

func (s *Service) GetMovie(ctx context.Context, req *GetMovieRequest) (*GetMovieResponse, error) {
	movie, err := s.movieRepository.GetMovie(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	genre := ""
	if len(movie.Genres) > 0 {
		genre = movie.Genres[0].Name
	}

	contract := &GetMovieResponse{}

	if movie == nil {
		m, err := s.movieAPIProxy.GetMovieData(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		if m == nil {
			return contract, nil
		}

		contract.Genre = genre
		contract.Title = m.Title
		contract.Year = m.Year
		contract.ID = m.ID

		mo := &models.Movie{
			ID:     m.ID,
			Title:  m.Title,
			Year:   m.Year,
			Genres: m.Genres,
		}

		if err := s.movieRepository.InsertMovie(ctx, mo); err != nil {
			return nil, err
		}

		return contract, nil
	}

	contract.Genre = genre
	contract.Title = movie.Title
	contract.Year = movie.Year
	contract.ID = movie.ID

	return contract, nil
}
