package env

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/justsushant/envbox/types"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(*http.Request) bool {
        return true
    },
}


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
	router.GET("/getTerminal/:id", h.getTerminal)
}


func (h *Handler) createEnv(c *gin.Context) {
	var payload types.CreateEnvPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"status": false, "error": err.Error()})
		return
	}

	resp, err := h.service.CreateEnv(h.client, payload)
	if err != nil {
		c.JSON(500, gin.H{"status": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "message": resp})
}

func (h *Handler) killEnv(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.service.KillEnv(h.client, id)
	if err != nil {
		c.JSON(500, gin.H{"status": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "message": resp})
}

func (h *Handler) getAllEnvs(c *gin.Context) {
	resp, err := h.service.GetAllEnvs()
	if err != nil {
		c.JSON(500, gin.H{"status": false, "error": err.Error()})
		return
	}
	if len(resp) == 0 {
		c.JSON(200, gin.H{"status": false, "error": "no envs found"})
		return	
	}

	c.JSON(200, gin.H{"status": true, "message": resp})
}

func (h *Handler) getTerminal(c *gin.Context) {
	// get the terminal hijaked response
	id := c.Param("id")
	termResp, err := h.service.GetTerminal(h.client, id)
	if err != nil {
		fmt.Println("Error while getting terminal: ", err)
		// c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer termResp.Close()

	// upgrade the connection to ws
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
		fmt.Println("Error while upgrading the connection: ", err)
        return
    }
    defer conn.Close()

	// reading from websocket
    go func() {
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                fmt.Println("read error:", err)
                return
            }
			// fmt.Println("ws: ", string(message))
			fmt.Fprint(termResp.Conn, string(message))
        }
    }()

	// Docker output to WebSocket loop
	for {
		buf := make([]byte, 1024)
		n, err := termResp.Reader.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		// fmt.Println("docker: ", string(buf[:n]))
		err = conn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

