package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GenerateSalt 生成随机盐值
func GenerateSalt(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// HashPasswordWithSalt 使用 HMAC-SHA256 + 盐值哈希密码，并进行 base64 编码
func HashPasswordWithSalt(password, salt string) string {
	// 使用 HMAC-SHA256 生成哈希
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(password))
	hashedBytes := h.Sum(nil)

	// 对哈希结果进行 base64 编码
	return base64.StdEncoding.EncodeToString(hashedBytes)
}
