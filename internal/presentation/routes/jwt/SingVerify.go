package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey string
	ttl       time.Duration
}

func NewJWTAdapter(secret string, ttl time.Duration) *JWT {
	return &JWT{secretKey: secret, ttl: ttl}
}

func (j *JWT) Sign(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	for k, v := range claims {
		tokenClaims[k] = v
	}
	tokenClaims["exp"] = time.Now().Add(j.ttl).Unix()

	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) Verify(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
