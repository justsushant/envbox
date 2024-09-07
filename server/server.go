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
	// router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins: true,
	// 	AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	// 	ExposeHeaders: []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return true
	// 	// },
	// 	MaxAge: 12 * time.Hour, // How long to cache the preflight response
	// }))


	router.LoadHTMLGlob("template/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"publicHost": "http://localhost:8080",
		})
	})

	router.GET("/terminal", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is required",
			})
			return
		}

		c.HTML(http.StatusOK, "terminal.html", gin.H{
			"publicHost": "localhost:8080",
			"id": id,
		})
	})

	

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
	envService := env.NewService(s.client, envStore, imageStore)
	envHandler := env.NewHandler(envService)
	envHandler.RegisterRoutes(apiRouter.Group("/env"))



	
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}