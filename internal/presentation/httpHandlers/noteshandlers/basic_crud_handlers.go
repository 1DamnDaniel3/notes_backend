package noteshandlers

import (
	"net/http"
	"notes_backend/internal/model"
	"notes_backend/internal/repository"
	ctxkeys "notes_backend/internal/repository/ctxKeys"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BasicNotesCrudHandlers struct {
	repo repository.INoteRepo
}

func NewBasicNotesCrudHandlers(repo repository.INoteRepo) *BasicNotesCrudHandlers {
	return &BasicNotesCrudHandlers{repo}
}

// -===================== METHODS =====================-

// GetOneNote godoc
// @Summary      GetOneNote
// @Description  GetOneNote
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "ID пользователя для поиска"
// @Success      200	{object} model.Note
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/notes/{id} [get]
func (h *BasicNotesCrudHandlers) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.repo.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uid, _ := strconv.Atoi(userIDFromContext)
	if uint(uid) != note.UserID { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own notes data"})
		return
	}

	resp := model.Note{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		Color:     note.Color,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// NoteCreate godoc
// @Summary      NoteCreate
// @Description  создание Note
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        input  body     model.Note  true  "Создание Note"
// @Success      200	{object} model.Note
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/notes [post]
func (h *BasicNotesCrudHandlers) Create(c *gin.Context) {
	ctx := c.Request.Context()
	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)

	note := model.Note{}
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, _ := strconv.Atoi(userIDFromContext)
	if note.UserID != uint(uid) { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only create your own notes"})
		return
	}

	if err := h.repo.Create(ctx, &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := model.Note{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		Color:     note.Color,
		IsPublic:  note.IsPublic,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	c.JSON(http.StatusCreated, output)
}

// NoteUpdate godoc
// @Summary      NoteUpdate
// @Description  обновление Note
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        input  body     model.Note  true  "обновление User"
// @Success      200	{object} map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/notes/{id} [patch]
func (h *BasicNotesCrudHandlers) Update(c *gin.Context) {
	ctx := c.Request.Context()

	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка

	note, err := h.repo.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uid, _ := strconv.Atoi(userIDFromContext)
	if note.UserID != uint(uid) { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only change your own notes data"})
		return
	}

	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// execute

	err = h.repo.Update(ctx, uint(id), fields)
	if err != nil {
		if err.Error() == "entity not found or access denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fields["id"] = id

	c.JSON(http.StatusOK, fields)
}

// NoteDelete godoc
// @Summary      NoteDelete
// @Description  удаление Note
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        id  path     int     true  "ID заметки для удаления"
// @Success      200	{object} DeleteResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/notes/{id} [delete]
func (h *BasicNotesCrudHandlers) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка

	note, err := h.repo.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uid, _ := strconv.Atoi(userIDFromContext)
	if uint(uid) != note.ID { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own notes"})
		return
	}

	// execute

	err = h.repo.Delete(ctx, uint(id), note)
	if err != nil {
		if err.Error() == "entity not found or access denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type DeleteResponse struct {
	ID int `json:"id"`
}

// GetAllUserNotes godoc
// @Summary      GetAllUserNotes
// @Description  Все заметки пользователя
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        public  query  bool  false  "Фильтр публичных заметок"
// @Success      200     {object} GetAllResponse
// @Failure      400     {object} map[string]interface{}
// @Failure      500     {object} map[string]interface{}
// @Router       /api/notes [get]
func (h *BasicNotesCrudHandlers) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var isPublic *bool

	isPublicStr := c.Query("public")
	if isPublicStr != "" {
		val, err := strconv.ParseBool(isPublicStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid public param"})
			return
		}
		isPublic = &val
	}

	notes, err := h.repo.GetAll(ctx, isPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := GetAllResponse{
		Data: make([]model.Note, len(*notes)),
	}

	copy(resp.Data, *notes)

	c.JSON(http.StatusOK, resp)
}

type GetAllResponse struct {
	Data []model.Note `json:"data"`
}
