package userhandlers

import (
	"net/http"
	"notes_backend/internal/model"
	netutils "notes_backend/internal/presentation/net_utils"
	"notes_backend/internal/service/jwt"
	"notes_backend/internal/service/userusecases"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc         userusecases.ILoginUC
	JwtService jwt.IJWT
}

func NewLoginHandler(uc userusecases.ILoginUC, JwtService jwt.IJWT) *AuthHandler {
	return &AuthHandler{uc, JwtService}
}

// Login godoc
// @Summary      Логин
// @Description  Вход по email/password.
// @Description  Web: JWT записывается в httpOnly Cookie.
// @Description  Mobile: JWT возвращается в body.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body     LoginDTO  true  "Данные для логина"
// @Param        X-Client-Type header string false "Тип клиента: mobile | web"
// @Success      200 {object} UserResponseDTO "Web response (JWT в cookie)"
// @Success      200 {object} LoginMobileResponse "Mobile response (JWT в body)"
// @Header       200 {string} Set-Cookie "JWT cookie (web only)"
// @Failure      400 {object} map[string]string
// @Router       /api/users/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	ctx := c.Request.Context()
	clientType := c.GetHeader("X-Client-Type") // mobile or web

	loginData := LoginDTO{}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:        loginData.Email,
		PasswordHash: loginData.Password,
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

	respBody := UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if clientType == "mobile" {
		respDto := LoginMobileResponse{
			User:  respBody,
			Token: token,
		}

		c.JSON(http.StatusOK, respDto)

	} else {

		c.SetCookie(
			"jwt",
			token,
			int(7*24*time.Hour.Seconds()),
			"/api",
			host,
			secure,
			true,
		)

		c.JSON(http.StatusOK, respBody)

	}

}

type LoginMobileResponse struct {
	User  UserResponseDTO `json:"user"`
	Token string          `json:"token"`
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

// AuthCheck godoc
// @Summary      Проверка авторизации
// @Description  Приходит кука с токеном, высылается 200, если токен ещё валиден
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 	{object}  AuthCheckResponse
// @Failure      401    {object}  map[string]string
// @Router       /api/users/auth-check [get]
func (a *AuthHandler) CheckAuth(c *gin.Context) {
	token, err := netutils.ExtractToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := a.JwtService.Verify(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// c.Set("user", claims)
	delete(claims, "exp")
	c.JSON(http.StatusOK, AuthCheckResponse{
		IsAuthenticated: true,
		User:            claims,
	})
}

type AuthCheckResponse struct {
	IsAuthenticated bool                   `json:"isAuthenticated"`
	User            map[string]interface{} `json:"user"`
}
