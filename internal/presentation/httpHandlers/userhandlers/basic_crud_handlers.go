package userhandlers

import (
	"net/http"
	"notes_backend/internal/model"
	"notes_backend/internal/repository"
	ctxkeys "notes_backend/internal/repository/ctxKeys"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BasicUserCrudHandlers struct {
	repo repository.IUserRepo
}

func NewBasicUserCrudHandlers(repo repository.IUserRepo) *BasicUserCrudHandlers {
	return &BasicUserCrudHandlers{repo}
}

// ====== DTO =======
type UserResponseDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// -===================== METHODS =====================-

// GetOneUser godoc
// @Summary      GetOneUser
// @Description  GetUserById
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "ID пользователя для поиска"
// @Success      200	{object} UserResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/users/{id} [get]
func (h *BasicUserCrudHandlers) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)
	idParam := c.Param("id")

	if userIDFromContext != idParam { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own data"})
		return
	}
	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// UserCreate godoc
// @Summary      UserCreate
// @Description  создание User
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        input  body     model.User  true  "Создание User"
// @Success      200	{object} UserResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/users [post]
func (h *BasicUserCrudHandlers) Create(c *gin.Context) {
	ctx := c.Request.Context()
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(user.PasswordHash) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long "})
		return
	}

	if err := h.repo.Create(ctx, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, output)
}

// UserUpdate godoc
// @Summary      UserUpdate
// @Description  обновление User
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        input  body     model.User  true  "обновление User"
// @Success      200	{object} map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/users/{id} [patch]
func (h *BasicUserCrudHandlers) Update(c *gin.Context) {
	ctx := c.Request.Context()

	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)

	idParam := c.Param("id")

	if userIDFromContext != idParam { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own data"})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// UserDelete godoc
// @Summary      UserDelete
// @Description  удаление User
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id  path     int     true  "ID пользователя для удаления"
// @Success      200	{object} DeleteResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/users/{id} [delete]
func (h *BasicUserCrudHandlers) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	userIDFromContext := ctx.Value(ctxkeys.UserId).(string)
	if userIDFromContext != idParam { // Проверка на вшивость
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
		return
	}
	id, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{}
	err = h.repo.Delete(ctx, uint(id), &user)
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
