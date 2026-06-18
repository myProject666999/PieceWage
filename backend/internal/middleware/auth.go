package middleware

import (
	"net/http"
	"strings"
	"time"

	"piece-wage/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	RealName string `json:"real_name"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64, username string, role int, realName string) (string, error) {
	cfg := config.AppConfig.JWT
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RealName: realName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "message": "未登录"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == auth {
			c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "message": "token格式错误"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "message": "token无效或已过期"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("realName", claims.RealName)
		c.Next()
	}
}

func RequireRole(roles ...int) gin.HandlerFunc {
	roleMap := make(map[int]bool, len(roles))
	for _, r := range roles {
		roleMap[r] = true
	}
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || !roleMap[role.(int)] {
			c.JSON(http.StatusForbidden, gin.H{"code": -1, "message": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
