package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/josevitorrodriguess/gochat/internal/api"
	"github.com/josevitorrodriguess/gochat/internal/repository"
	"github.com/josevitorrodriguess/gochat/internal/storage/postgres"
	"github.com/josevitorrodriguess/gochat/internal/usecase"
)

var (
	userRepo    repository.UserRepository
	userUseCase usecase.UserUseCase
)

func main() {
	api := api.NewApi(userUseCase)

	api.RegisterRoutes()

	log.Println("Servidor iniciando na porta :3050")
	log.Fatal(api.Start(":3050"))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	pgStorage, err := postgres.NewDefault()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db := pgStorage.GetDB()

	userRepo = repository.NewUserRepository(db)
	userUseCase = usecase.NewUserUseCase(userRepo)
}
