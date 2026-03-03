package routes

import (
	"notes_backend/internal/presentation/middleware"
	"notes_backend/internal/presentation/routes/jwt"
	userroutes "notes_backend/internal/presentation/routes/user_routes"
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
	JWTSigner := jwt.NewJWTAdapter(secret, 5*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(*JWTSigner)

	protected := api
	if authMiddleware != nil {
		protected = r.Group("")
		protected.Use(authMiddleware.TryAuth())
	}

	userroutes.UserBasicRoutes(api, protected, db, authMiddleware)
}
