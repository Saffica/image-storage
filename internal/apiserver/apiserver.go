package apiserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	config *Config
	router *gin.Engine
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: gin.Default(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureRouter(); err != nil {
		return err
	}

	return nil
}

func (s *APIServer) configureRouter() error {
	s.router.GET("/img", func(ctx *gin.Context) {
		hash := ctx.Query("hash")
		message := fmt.Sprintf("hash: %s", hash)
		ctx.String(http.StatusOK, message)
	})

	if err := s.router.Run(s.config.BindAddr); err != nil {
		return err
	}

	return nil
}
