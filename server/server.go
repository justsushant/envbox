package server

import (
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
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
	
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}