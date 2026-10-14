package main

import (
	"auth/database"
	"auth/handler"
	"auth/repository"
	"auth/service"

	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbURL := os.Getenv("DB_URL")
	conn, err := database.ConnectPostgres(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	userRepository := repository.NewPostrgresUserRepository(conn)

	userService := service.NewUserService(userRepository)

	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/register", userHandler.Register)
	mux.HandleFunc("POST /users/login", userHandler.Login)

	fmt.Println("Server running on port: localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
