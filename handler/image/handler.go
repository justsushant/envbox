package image

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/envbox/types"
)
type Handler struct {
	service types.ImageService
}

func NewHandler(service types.ImageService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/getImages", h.GetImages)
}

func (h *Handler) GetImages(c *gin.Context) {
	data, err := h.service.GetImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "no images found",
		})
		return
	}
	c.JSON(http.StatusOK, data)
}