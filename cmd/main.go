package main

import (
	"avito/pkg/handler"
	"avito/pkg/repository"
	"avito/pkg/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

// @title           Avito Test Assignment
// @version         1.0

// @contact.name   GADJIIAVOV DJAMAL
// @contact.url    https://dj.ama1.ru
// @contact.email  mail@dj.ama1.ru

// @host      localhost:8000
// @BasePath  /api/

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("[Config Error] Initialization error: %s", err.Error())
	}

	db, err := repository.NewDB(repository.Config{
		Driver:   "postgres",
		Name:     viper.GetString("POSTGRES_DB"),
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		Username: viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
	})

	if err != nil {
		log.Fatalf("[DB ERROR] %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(Server)
	go func() {
		if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("[Http Server Error] %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("[Http Server Error] %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("[DB Error] %s", err.Error())
	}
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func initConfig() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}
