package jwt

import (
	"kama_chat_server/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var Middleware = func(c *gin.Context) {
	// 从 Authorization 头获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		helper.JsonBack(c, "登录信息已失效", -1, "Missing Authorization header")
		return
	}

	// 通常格式为 "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		helper.JsonBack(c, "登录信息已失效", -1, "Invalid Authorization header format")
		return
	}
	tokenString := authHeader[len(bearerPrefix):]

	// 解析 token
	claims, err := parseJWT(tokenString)
	if err != nil {
		helper.JsonBack(c, "登录信息已失效", -1, "Invalid or expired token")
		return
	}

	// 将用户信息存入上下文，供后续处理使用
	c.Set("user_id", claims.UserId)
	c.Next()
}

// 用于签名 JWT 的密钥（实际应用中应从配置文件或环境变量读取）
var jwtSecret = []byte("it'smesharve")

// 用户登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 自定义 JWT 的 claims
type UserClaims struct {
	UserId string `json:"username"`
	jwt.RegisteredClaims
}

// 生成 JWT
func generateJWT(userId string) (string, error) {
	claims := UserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 有效期24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sharve",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析并验证 JWT
func parseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// func main() {
// 	r := gin.Default()

// 	// 登录接口，返回 JWT
// 	r.POST("/login", func(c *gin.Context) {
// 		var req LoginRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// 硬编码验证（实际应用中应查询数据库并验证密码哈希）
// 		if req.Username == "admin" && req.Password == "123456" {
// 			token, err := generateJWT(req.Username)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 				return
// 			}
// 			c.JSON(http.StatusOK, gin.H{"token": token})
// 			return
// 		}

// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 	})

// 	// 受保护的路由组，使用 authMiddleware
// 	authorized := r.Group("/api")
// 	authorized.Use(authMiddleware())
// 	{
// 		authorized.GET("/user", func(c *gin.Context) {
// 			username, _ := c.Get("username")
// 			c.JSON(http.StatusOK, gin.H{
// 				"message":  "Welcome to the protected resource!",
// 				"username": username,
// 			})
// 		})
// 	}

// 	// 启动服务
// 	r.Run(":8080")
// }
