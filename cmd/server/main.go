package main

import (
	"log"
	"notes_backend/internal/model"
	"notes_backend/internal/presentation/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// swag init -g cmd/server/main.go

// @title Notes API
// @version 1.0
// @description Simple notes service
// @host localhost:3001
// @BasePath /
func main() {
	_ = godotenv.Load()
	PORT := os.Getenv("PORT")
	router := gin.Default()

	if os.Getenv("ENV") == "dev" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}))
	}

	db, err := model.ConnectDB()
	if err != nil {
		log.Fatalf("failed to load config :( because: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Note{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	log.Println("Database migrated successfully")

	routes.SetupRoutes(router, db) // routes

	router.Run(":" + PORT)

}
