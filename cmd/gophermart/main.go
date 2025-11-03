package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lovchik/gophermart/internal/server/config"
	"github.com/Lovchik/gophermart/internal/server/feign"
	"github.com/Lovchik/gophermart/internal/server/handlers"
	"github.com/Lovchik/gophermart/internal/server/storage"
	"github.com/Lovchik/gophermart/internal/server/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

func main() {
	Serve()
}

func Serve() {
	s := &handlers.Service{}
	s.Config = config.InitConfig()
	s.Validator = validator.New(validator.WithRequiredStructEnabled())
	s.JwtKeysPair = utils.InitJwtPair(s.Config.PrivateKey, s.Config.PublicKey)
	accuralFeign, err := feign.NewAccuralFeign(s.Config.AccuralSystemAddress)
	if err != nil {
		log.Fatalf("Failed to connect to accural system: %v", err)
	}
	s.Feign = accuralFeign
	s.WebServer = gin.Default()
	pgStore, err := storage.NewPgStorage(context.Background(), s.Config.DatabaseDNS)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	s.Store = pgStore
	s.WebServer.Use(CORSMiddleware(s.Config.AllowOrigins))
	api := s.WebServer.Group("/api")
	handlers.UserRegister(api.Group("user"), s)

	httpServer := &http.Server{Addr: s.Config.Address, Handler: s.WebServer}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("http server shutdown error: %v", err)
	}
	if s.Store != nil {
		if err := s.Store.Close(ctx); err != nil {
			log.Errorf("postgres close error: %v", err)
		}
	}
	log.Info("server exited")
}

func CORSMiddleware(allowOrigins string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: strings.Split(allowOrigins, ","),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "Access-Control-Allow-Headers", "Content-Type", "access-control-allow-origin", "access-control-allow-headers"},
		MaxAge:       12 * time.Hour,
	})
}
