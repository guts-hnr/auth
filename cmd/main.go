package main

import (
	"auth/database"
	"auth/handler"
	"auth/middleware"
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

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	dbURL := os.Getenv("DB_URL")
	conn, err := database.ConnectPostgres(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	userRepository := repository.NewPostrgresUserRepository(conn)

	userService := service.NewUserService(userRepository, jwtSecret)

	userHandler := handler.NewUserHandler(userService)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle("POST /profile", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.Profile)))

	fmt.Println("Server running on port: localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
