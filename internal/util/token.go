package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken 计算 Token 的哈希值（SHA256）
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
