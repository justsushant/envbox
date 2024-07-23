package server

import (
	"database/sql"
	"net/http"
	"log"
	"github.com/gin-gonic/gin"

	"github.com/justsushant/envbox/handler/image"
)

type Server struct {
	addr string
	db *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db: db,
	}
}

func (s *Server) Run() error {
	router := gin.Default()
	subRouter := router.Group("/api/v1")

	subRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	imageStore := image.NewStore(s.db)
	imageService := image.NewService(imageStore)
	imageHandler := image.NewHandler(imageService)
	imageHandler.RegisterRoutes(subRouter.Group("/image"))
	
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}