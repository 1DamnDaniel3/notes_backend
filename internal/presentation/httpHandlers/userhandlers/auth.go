package userhandlers

import (
	"net/http"
	"notes_backend/internal/model"
	"notes_backend/internal/service/userusecases"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc userusecases.ILoginUC
}

func NewLoginHandler(uc userusecases.ILoginUC) *AuthHandler {
	return &AuthHandler{uc}
}

// Login godoc
// @Summary      Логин
// @Description  Вход стандарт email password, запись в httpOnly Cookies JWT, в body лежит user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body     LoginDTO  true  "Данные для логина"
// @Success      200	{object} UserResponseDTO
// @Header 		 200	{string} Set-Cookie "JWT-токен"
// @Failure      400    {object}  map[string]string
// @Router       /api/users/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	ctx := c.Request.Context()
	loginData := LoginDTO{}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		// ID:       0,
		// Nickname: "",

		Email:        loginData.Email,
		PasswordHash: loginData.Password,

		// CreatedAt: time.Time{},
		// UpdatedAt: time.Time{},
	}

	token, err := h.uc.Execute(ctx, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secure := os.Getenv("ENV") == "prod"
	var host string
	if os.Getenv("ENV") == "prod" {
		host = os.Getenv("HOST")
	} else {
		host = "localhost"
	}

	c.SetCookie(
		"jwt",
		token,
		int(5*time.Hour.Seconds()),
		"/api",
		host,
		secure,
		true,
	)

	respBody := UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, respBody)
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Logout godoc
// @Summary      Логаут
// @Description  Простой get для logout. Пока только удаляет JWT из браузера, не ведёт blacklist
// @Tags         Auth
// @Success      204
// @Router       /api/users/logout [get]
func (r *AuthHandler) Logout(c *gin.Context) {

	secure := os.Getenv("ENV") == "prod"
	var host string
	if os.Getenv("ENV") == "prod" {
		host = os.Getenv("HOST")
	} else {
		host = "localhost"
	}

	c.SetCookie(
		"jwt",
		"",
		-1,
		"/api",
		host,
		secure,
		true,
	)
	c.Status(http.StatusNoContent)
}
