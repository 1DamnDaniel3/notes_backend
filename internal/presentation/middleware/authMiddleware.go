package middleware

import (
	"context"
	"net/http"
	netutils "notes_backend/internal/presentation/net_utils"
	ctxkeys "notes_backend/internal/repository/ctxKeys"
	"notes_backend/internal/service/jwt"

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
		token, err := netutils.ExtractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
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

		// log.Println("=================", claims["user_id"].(string))

		user_id, ok := claims["user_id"].(float64)
		if !ok {
			c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in claims"})
			return
		}

		userID := uint(user_id)

		ctx := context.WithValue(
			c.Request.Context(),
			ctxkeys.UserId,
			userID,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Set("user", claims)
		c.Next()
	}
}
