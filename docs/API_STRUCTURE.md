# API 结构体组织说明

## 📁 目录结构

```
internal/controller/
├── request/          # 请求结构体
│   ├── api.go       # API 请求结构体
│   ├── admin.go     # 管理请求结构体
│   └── example.go   # 示例请求结构体
├── response/         # 响应结构体
│   ├── common.go    # 通用响应结构体
│   ├── api.go       # API 响应结构体
│   ├── admin.go     # 管理响应结构体
│   └── example.go   # 示例响应结构体
├── api.go           # API 控制器
├── admin.go         # 管理控制器
└── web.go           # Web 控制器
```

## 🏗️ 结构体设计原则

### 1. 请求结构体 (Request)

- 使用 `binding` 标签进行参数验证
- 支持 JSON、Form、Query 等不同绑定方式
- 命名规范：`{功能}Request`

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

### 2. 响应结构体 (Response)

- 继承 `BaseResponse` 基础结构
- 统一的错误码和消息格式
- 命名规范：`{功能}Response`

```go
type LoginResponse struct {
    BaseResponse
    Token string `json:"token,omitempty"`
}
```

### 3. 基础响应结构 (BaseResponse)

```go
type BaseResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

## 📝 使用示例

### 1. 定义请求结构体

```go
// internal/controller/request/user.go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### 2. 定义响应结构体

```go
// internal/controller/response/user.go
type CreateUserResponse struct {
    BaseResponse
    Data UserData `json:"data"`
}

type UserData struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}
```

### 3. 在控制器中使用

```go
// internal/controller/user.go
func (u *UserController) CreateUser(c *gin.Context) {
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "参数错误"))
        return
    }
    
    // 业务逻辑处理
    user := u.userService.CreateUser(req)
    
    c.JSON(http.StatusOK, response.CreateUserResponse{
        BaseResponse: response.BaseResponse{
            Code:    0,
            Message: "创建成功",
        },
        Data: user,
    })
}
```

## 🔧 验证规则

### 常用验证标签

- `required`: 必填字段
- `min=n`: 最小长度/值
- `max=n`: 最大长度/值
- `oneof=value1 value2`: 枚举值
- `omitempty`: 空值时不验证

### 示例

```go
type UserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Age      int    `json:"age" binding:"min=1,max=120"`
    Status   string `json:"status" binding:"oneof=active inactive"`
}
```

## 📊 响应码规范

### 成功响应
- `0`: 成功

### 客户端错误
- `400`: 请求参数错误
- `401`: 未授权
- `403`: 禁止访问
- `404`: 资源不存在

### 服务器错误
- `500`: 内部服务器错误

## 🎯 最佳实践

1. **命名规范**: 使用清晰、描述性的名称
2. **结构体复用**: 通过组合复用基础结构
3. **验证完整**: 充分利用 Gin 的验证功能
4. **文档注释**: 为每个结构体添加注释
5. **类型安全**: 使用强类型而非 `any`

## 📚 扩展指南

### 添加新的 API

1. 在 `request/` 目录下定义请求结构体
2. 在 `response/` 目录下定义响应结构体
3. 在控制器中实现处理逻辑
4. 注册路由

### 添加新的验证规则

1. 使用 Gin 内置验证器
2. 自定义验证器（如需要）
3. 在结构体标签中应用

### 添加新的响应类型

1. 继承 `BaseResponse`
2. 添加特定字段
3. 提供构造函数（如需要）
