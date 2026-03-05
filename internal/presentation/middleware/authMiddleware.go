package middleware

import (
	"context"
	"net/http"
	ctxkeys "notes_backend/internal/repository/ctxKeys"
	"notes_backend/internal/service/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	JwtService jwt.IJWT
}

func NewAuthMiddleware(JwtService jwt.IJWT) *AuthMiddleware {
	return &AuthMiddleware{JwtService}
}

func (a *AuthMiddleware) TryAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt")
		if err != nil {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := a.JwtService.Verify(token)
		if err != nil {
			c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		user_id, ok := claims["user_id"].(string)
		if !ok {
			c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in claims"})
			return
		}

		ctx := context.WithValue(
			c.Request.Context(),
			ctxkeys.UserId,
			user_id,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Set("user", claims)
		c.Next()
	}
}
