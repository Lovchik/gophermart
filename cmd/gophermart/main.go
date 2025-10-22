package main

import (
	"context"
	"github.com/Lovchik/gophermart/internal/server/config"
	"github.com/Lovchik/gophermart/internal/server/handlers"
	"github.com/Lovchik/gophermart/internal/server/models"
	"github.com/Lovchik/gophermart/internal/server/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
