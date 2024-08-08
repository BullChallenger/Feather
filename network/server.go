package network

import (
	"feather/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	engine  *gin.Engine
	service *service.Service
	port    string
	ip      string
}

func NewServer(service *service.Service, port string) *Server {
	s := &Server{engine: gin.New(), service: service, port: port}
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	registerServer(s)
	return s
}

func (s *Server) StartServer() error {
	log.Println("Go Server Starting...")
	return s.engine.Run(s.port)
}
