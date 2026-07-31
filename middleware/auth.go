package middleware

import (
	"errors"
	"net/http"

	"github.com/MeiSiristhebest/tiktok-backend-go/config"
	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "tiktok-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.SecretKey))
}

func ParseToken(tokenStr string) (int64, error) {
	if tokenStr == "" {
		return 0, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return claims.UserID, nil
}

// AuthMiddleware 强强制鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenStr = authHeader[7:]
			}
		}

		userID, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, model.BaseResponse{
				StatusCode: 401,
				StatusMsg:  "User authentication failed: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("current_user_id", userID)
		c.Next()
	}
}

// OptionalAuthMiddleware 可选鉴权中间件 (如 Feed 接口)
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		if tokenStr != "" {
			userID, err := ParseToken(tokenStr)
			if err == nil {
				c.Set("current_user_id", userID)
			}
		}
		c.Next()
	}
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
