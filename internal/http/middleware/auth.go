package middleware

import (
  "errors"
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"

  "DCS/internal/auth"
)

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
      })
      return
    }

    claims, err := jwtManager.Verify(tokenString)
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
