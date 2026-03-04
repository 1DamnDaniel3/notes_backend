package routes

import (
	"notes_backend/internal/presentation/middleware"
	userroutes "notes_backend/internal/presentation/routes/user_routes"
	hashservice "notes_backend/internal/service/hashService"
	"notes_backend/internal/service/jwt"
	"os"
	"time"

	_ "notes_backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// == JWT ==
	secret := os.Getenv("JWT_SECRET")
	jwtService := jwt.NewJWTAdapter(secret, 5*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// == hash ==
	hashService := hashservice.NewbcryptHashService()

	protected := api
	if authMiddleware != nil {
		protected = api.Group("")
		protected.Use(authMiddleware.TryAuth())
	}

	userroutes.UserRoutes(api, protected, db, authMiddleware, hashService, jwtService)
}
