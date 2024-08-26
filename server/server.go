package server

import (
	"database/sql"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/docker/docker/client"

	"github.com/justsushant/envbox/handler/image"
	"github.com/justsushant/envbox/handler/env"
)

type Server struct {
	addr string
	db *sql.DB
	client *client.Client
}

func NewServer(addr string, db *sql.DB, client *client.Client) *Server {
	return &Server{
		addr: addr,
		db: db,
		client: client,
	}
}

func (s *Server) Run() error {
	router := gin.Default()

	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	imageStore := image.NewStore(s.db)
	imageService := image.NewService(imageStore)
	imageHandler := image.NewHandler(imageService)
	imageHandler.RegisterRoutes(apiRouter.Group("/image"))

	envStore := env.NewStore(s.db)
	envService := env.NewService(envStore, imageStore)
	envHandler := env.NewHandler(envService, s.client)
	envHandler.RegisterRoutes(apiRouter.Group("/env"))



	
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}