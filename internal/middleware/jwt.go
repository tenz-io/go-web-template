package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
)

// JWT Claims 结构体
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	Iss      string `json:"iss"`
	Sub      string `json:"sub"`
}

// JWT 管理器
type JWTManager struct {
	secret     []byte
	expireTime time.Duration
	issuer     string
}

// 创建 JWT 管理器
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:     []byte(cfg.Secret),
		expireTime: cfg.ExpireTime,
		issuer:     cfg.Issuer,
	}
}

// 生成 JWT Token
func (j *JWTManager) GenerateToken(userID, username, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Exp:      now.Add(j.expireTime).Unix(),
		Iat:      now.Unix(),
		Iss:      j.issuer,
		Sub:      userID,
	}

	// 创建 JWT header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 编码 header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// 编码 payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := j.createSignature(message)

	// 组合 JWT
	return message + "." + signature, nil
}

// 验证 JWT Token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// 验证签名
	message := parts[0] + "." + parts[1]
	expectedSignature := j.createSignature(message)
	if parts[2] != expectedSignature {
		return nil, errors.New("invalid signature")
	}

	// 解码 payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, err
	}

	// 检查过期时间
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// 创建 HMAC-SHA256 签名
func (j *JWTManager) createSignature(message string) string {
	h := hmac.New(sha256.New, j.secret)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// 鉴权类型
type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeBearer AuthType = "bearer"
	AuthTypeCookie AuthType = "cookie"
)

// 鉴权配置
type AuthConfig struct {
	Type     AuthType `json:"type"`
	Required bool     `json:"required"`
}

// 鉴权中间件
func Auth(config AuthConfig, jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Required {
			c.Next()
			return
		}

		switch config.Type {
		case AuthTypeNone:
			c.Next()
		case AuthTypeBearer:
			handleBearerAuth(c, jwtManager)
		case AuthTypeCookie:
			handleCookieAuth(c, jwtManager)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "不支持的鉴权类型",
			})
			c.Abort()
		}
	}
}

// Bearer Token 认证
func handleBearerAuth(c *gin.Context, jwtManager *JWTManager) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "缺少 Authorization 头",
		})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的 Bearer token",
		})
		c.Abort()
		return
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		logger.FromContext(c.Request.Context()).Warn("JWT validation failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的 token",
		})
		c.Abort()
		return
	}

	// 设置用户信息到上下文
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
	c.Next()
}

// Cookie 认证
func handleCookieAuth(c *gin.Context, jwtManager *JWTManager) {
	token, err := c.Cookie("jwt_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录",
		})
		c.Abort()
		return
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		logger.FromContext(c.Request.Context()).Warn("JWT validation failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的 session",
		})
		c.Abort()
		return
	}

	// 设置用户信息到上下文
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
	c.Next()
}

// 管理员权限中间件
func AdminAuth(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt_token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.FromContext(c.Request.Context()).Warn("Admin JWT validation failed")
			c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
			c.Abort()
			return
		}

		// 检查是否为管理员角色
		if claims.Role != "admin" {
			c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
