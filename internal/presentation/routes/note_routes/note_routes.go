package noteroutes

import (
	"notes_backend/internal/presentation/httpHandlers/noteshandlers"
	"notes_backend/internal/presentation/middleware"
	"notes_backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NoteRoutes(
	r *gin.RouterGroup,
	protected *gin.RouterGroup,
	db *gorm.DB,
	authMiddleware *middleware.AuthMiddleware,

) {
	noteRepo := repository.NewNoteRepo(db)
	basicHandlers := noteshandlers.NewBasicNotesCrudHandlers(noteRepo)

	noteQueryService := repository.NewNoteQueryService(db)
	queryServiceHandler := noteshandlers.NewQueryServiceHandler(noteQueryService)

	r.GET("/notes/public/:page", queryServiceHandler.GetAllPublic) // GetAllPublic

	// ---======= protected routes
	protected.POST("/notes", basicHandlers.Create)       // Create
	protected.GET("/notes", basicHandlers.GetAll)        // GetAll
	protected.PATCH("/notes/:id", basicHandlers.Update)  // Update
	protected.GET("/notes/:id", basicHandlers.GetByID)   // GetById
	protected.DELETE("/notes/:id", basicHandlers.Delete) // Delete

}
