package userroutes

import (
	"notes_backend/internal/presentation/httpHandlers/userhandlers"
	"notes_backend/internal/presentation/middleware"
	"notes_backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserBasicRoutes(
	r *gin.RouterGroup,
	protected *gin.RouterGroup,
	db *gorm.DB,
	authMiddleware *middleware.AuthMiddleware,
) {
	userRepo := repository.NewUserRepo(db)
	basicHandlers := userhandlers.NewBasicUserCrudHandlers(userRepo)

	protected.POST("users", basicHandlers.Create) // Create
}
