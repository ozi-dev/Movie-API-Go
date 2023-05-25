package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"selfuelAPI/internal/service"
)

type Server struct {
	srv *service.Service
}

func New(srv *service.Service) *Server {
	return &Server{
		srv: srv,
	}
}

func (s *Server) GetMovie(c echo.Context) error {
	req := &service.GetMovieRequest{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	movie, err := s.srv.GetMovie(c.Request().Context(), req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, movie)
}
