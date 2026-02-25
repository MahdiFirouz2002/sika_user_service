package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sika/internal/database"
	"sika/internal/repositoy"
	"sika/internal/service"
	"sika/pkg/config"
	"syscall"
	"time"

	"sika/internal/server/controller"
	httpServer "sika/internal/server/http"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		godotenv.Load()
	}
}

func main() {
	mode := flag.String("mode", "server", "chose mode btwen usage and api handler")
	flag.Parse()

	// database
	cfg := config.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName, cfg.DatabasePort)

	database, err := database.NewConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "server":
		runServer(*database, cfg.HTTPServerHost, cfg.HTTPServerPort)
	case "insert":
		runInsert(*database)
	default:
		log.Println("mode flag required")
	}
}

func runServer(database database.Database, host, port string) {
	userRepository := repositoy.NewUserRepositor(database)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	server := httpServer.NewHttpServer(host, port, *userController)

	go func() {
		log.Fatal(server.Listen())
	}()

	<-ctx.Done()
	log.Println("shuting down server gracefully")
	server.ShutDown()
}

func runInsert(database database.Database) {
	userRepository := repositoy.NewUserRepositor(database)
	userService := service.NewUserService(userRepository)
	producer := service.NewProducer(userService, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	producer.RunInsert(ctx)
}
