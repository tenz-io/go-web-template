package constant

// 用户角色常量
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// 角色列表
var AllRoles = []Role{
	RoleUser,
	RoleAdmin,
}

// IsValidRole 验证角色是否有效
func IsValidRole(role string) bool {
	for _, r := range AllRoles {
		if string(r) == role {
			return true
		}
	}
	return false
}
