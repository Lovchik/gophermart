package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gofermart/internal/server/config"
	"gofermart/internal/server/handlers"
	"gofermart/internal/server/models"
	"gofermart/internal/server/storage"
)

func main() {
	config.InitConfig()
	models.InitJwtPair()
	Serve()
}

func Serve() {
	s := &handlers.Service{}
	s.WebServer = gin.Default()
	pgStore, err := storage.NewPgStorage(context.Background(), config.GetConfig().DatabaseDNS)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	s.Store = pgStore
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowAllOrigins = true
	s.WebServer.Use(cors.New(ginConfig))
	api := s.WebServer.Group("/api")
	handlers.UserRegister(api.Group("user"), s)
	err = s.WebServer.Run(config.GetConfig().Address)
	if err != nil {
		log.Error(err)
		return
	}
}
