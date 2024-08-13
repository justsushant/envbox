package env

import (
	"encoding/json"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/justsushant/envbox/types"
)

type Handler struct {
	service types.EnvService
	client  *client.Client
}

func NewHandler(service types.EnvService, client *client.Client) *Handler {
	return &Handler{
		service: service,
		client:  client,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/getAllEnvs", h.getAllEnvs)
	router.POST("/createEnv", h.createEnv)
	router.PATCH("/killEnv/:id", h.killEnv)
}

func (h *Handler) createEnv(c *gin.Context) {
	var payload types.CreateEnvPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateEnv(h.client, payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": resp})
}

func (h *Handler) killEnv(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.service.KillEnv(h.client, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": resp})
}

func (h *Handler) getAllEnvs(c *gin.Context) {
	resp, err := h.service.GetAllEnvs()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": resp})
}
