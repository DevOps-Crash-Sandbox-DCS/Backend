package middleware

import (
<<<<<<< HEAD
=======
<<<<<<< HEAD
	"errors"
=======
>>>>>>> feature/ai-hint-service
>>>>>>> front_and_ai_service
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"DCS/internal/auth"
)

<<<<<<< HEAD
=======
<<<<<<< HEAD
var (
	errAuthorizationRequired   = errors.New("authorization header is required")
	errInvalidAuthHeaderFormat = errors.New("invalid authorization header format")
)

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
=======
>>>>>>> front_and_ai_service
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
<<<<<<< HEAD
=======
>>>>>>> feature/ai-hint-service
>>>>>>> front_and_ai_service
			})
			return
		}

<<<<<<< HEAD
=======
<<<<<<< HEAD
		claims, err := jwtManager.Verify(tokenString)
=======
>>>>>>> front_and_ai_service
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		claims, err := jwtManager.Verify(parts[1])
<<<<<<< HEAD
=======
>>>>>>> feature/ai-hint-service
>>>>>>> front_and_ai_service
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token payload",
			})
			return
		}

		c.Set("userId", userID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
<<<<<<< HEAD
=======
<<<<<<< HEAD

// extractToken берёт JWT из заголовка Authorization: Bearer <token>,
// а если заголовка нет — из query-параметра ?token=... Второй вариант нужен
// для WebSocket-подключений (/terminal): браузер не умеет отправлять
// произвольные заголовки при WebSocket-хендшейке, поэтому фронт передаёт
// токен прямо в адресе.
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", errInvalidAuthHeaderFormat
		}
		return parts[1], nil
	}

	if token := c.Query("token"); token != "" {
		return token, nil
	}

	return "", errAuthorizationRequired
}
=======
>>>>>>> feature/ai-hint-service
>>>>>>> front_and_ai_service
