package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go-web-template/internal/constant"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
)

const (
	// JWTTokenCookieName JWT token cookie 名称
	JWTTokenCookieName = "jwt_token"
	userIdName         = "user_id"
	roleName           = "role"
)

// JWT Claims 结构体
type Claims struct {
	UserID int64 `json:"user_id"`
	Role   int32 `json:"role"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
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

func (j *JWTManager) GenerateToken(userID int64, role int32) (string, error) {
	return j.GenerateTokenWithExpire(userID, role, j.expireTime)
}

// 生成 JWT Token
func (j *JWTManager) GenerateTokenWithExpire(userID int64, role int32, expDuration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Exp:    now.Add(expDuration).Unix(),
		Iat:    now.Unix(),
	}

	// 创建 JWT header
	header := map[string]any{
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
	Type     AuthType      `json:"type"`
	Required bool          `json:"required"`
	Role     constant.Role `json:"role"`
}

// 鉴权中间件
func Auth(config AuthConfig, jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		le := logger.FromContext(c.Request.Context())

		if !config.Required {
			le.Debug("authentication not required, skipping")
			c.Next()
			return
		}

		switch config.Type {
		case AuthTypeNone:
			le.Debug("configured to skip authentication")
			c.Next()
		case AuthTypeBearer:
			handleBearerAuth(c, jwtManager)
		case AuthTypeCookie:
			handleCookieAuth(c, jwtManager, config.Role)
		default:
			le.Warn("unknown auth type, aborting")
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
	le := logger.FromContext(c.Request.Context())

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		le.Warn("Authorization header is empty")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "缺少 Authorization 头",
		})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		le.Warn("Authorization header is valid")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的 Bearer token",
		})
		c.Abort()
		return
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		le.Warn("JWT validation failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的 token",
		})
		c.Abort()
		return
	}

	le.Debug("validate token")

	// 设置用户信息到上下文
	c.Set(userIdName, claims.UserID)
	c.Set(roleName, claims.Role)

	c.Next()
}

func GetUserInfoFromContext(c *gin.Context) (userID int64, role int32, err error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		err = errors.New("user_id not found in context")
		return
	}

	roleInterface, exists := c.Get("role")
	if !exists {
		err = errors.New("role not found in context")
		return
	}

	var ok bool
	if userID, ok = userIDInterface.(int64); !ok {
		err = errors.New("invalid user_id type")
		return
	}

	role, ok = roleInterface.(int32)
	if !ok {
		err = errors.New("invalid role type")
		return
	}

	return userID, role, nil
}

// Cookie 认证（通用）
func handleCookieAuth(c *gin.Context, jwtManager *JWTManager, role constant.Role) {
	le := logger.FromContext(c.Request.Context())

	token, err := c.Cookie(JWTTokenCookieName)
	if err != nil {
		le.Warn("JWT cookie auth failed")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		le.Warn("JWT validation failed")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}

	// 检查角色权限

	switch role {
	case constant.RoleAdmin:
		// 管理员权限：必须是管理员角色
		if claims.Role != int32(constant.RoleAdmin) {
			le.Warn("JWT role is invalid")
			logger.FromContext(c.Request.Context()).Warn("Admin authentication failed: insufficient role")
			c.Redirect(http.StatusTemporaryRedirect, "/login?error=权限不足")
			c.Abort()
			return
		}
	default:
		//ignore
	}

	le.Debug("validate token")

	// 设置用户信息到上下文
	c.Set(userIdName, claims.UserID)
	c.Set(roleName, claims.Role)
	c.Next()
}
