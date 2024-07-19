package main

import (
	"fmt"
	"github.com/justsushant/envbox/server"
	"github.com/justsushant/envbox/config"
)

func main() {
	s := server.NewServer(fmt.Sprintf(":%s", config.Envs.Port))
	s.Run()
}