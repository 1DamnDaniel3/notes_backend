package userroutes

import (
	"notes_backend/internal/presentation/httpHandlers/userhandlers"
	"notes_backend/internal/presentation/middleware"
	"notes_backend/internal/repository"
	hashservice "notes_backend/internal/service/hashService"
	"notes_backend/internal/service/jwt"
	"notes_backend/internal/service/userusecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(
	r *gin.RouterGroup,
	protected *gin.RouterGroup,
	db *gorm.DB,
	authMiddleware *middleware.AuthMiddleware,

	hashService hashservice.IBcryptHashService,
	jwtService jwt.IJWT,
) {
	userRepo := repository.NewUserRepo(db)
	basicHandlers := userhandlers.NewBasicUserCrudHandlers(userRepo)

	// auth
	loginUC := userusecases.NewLoginUC(userRepo, hashService, jwtService)
	authHandler := userhandlers.NewLoginHandler(loginUC, jwtService)

	r.POST("/users", basicHandlers.Create)     // Create
	r.POST("/users/login", authHandler.Login)  // Login
	r.GET("/users/logout", authHandler.Logout) // Logout

	// ---======= protected routes

	protected.PATCH("/users", basicHandlers.Update)  // Update
	protected.GET("/users", basicHandlers.GetByID)   // GetById
	protected.DELETE("/users", basicHandlers.Delete) // Delete

}
