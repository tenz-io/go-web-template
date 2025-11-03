# 🚀 Go Web 模板项目

一个现代化的 Go Web 应用模板，提供完整的 API 接口和管理后台功能。

## ✨ 特性

- 🔧 **RESTful API** - 完整的 API 接口，支持 JSON 格式数据交换
- 🛡️ **权限管理** - 内置 JWT 认证和基于角色的权限控制系统
- 🎛️ **管理后台** - 现代化的管理界面，支持用户管理和系统监控
- 📊 **监控日志** - 完整的日志记录和系统监控功能
- 🏗️ **依赖注入** - 使用 Google Wire 进行依赖注入
- 🧪 **测试支持** - 内置 Mock 生成和测试框架
- 🐳 **Docker 支持** - 提供 Docker 配置和部署脚本

## 🏗️ 项目结构

```
go-web-template/
├── api/                    # API 定义和生成的文件
│   ├── http/app/          # HTTP API 定义
│   └── custom/            # 自定义 protobuf 定义
├── bin/                   # 编译后的二进制文件
├── cmd/                   # 应用程序入口
│   ├── main.go           # 主程序入口
│   └── webgo/            # Web 服务器
├── config/                # 配置文件
│   └── app.yaml          # 应用配置
├── internal/              # 内部包
│   ├── config/           # 配置结构
│   ├── controller/       # 控制器层
│   │   ├── request/      # 请求结构体
│   │   └── response/     # 响应结构体
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   ├── setup/            # 依赖注入配置
│   └── util/             # 工具函数
├── web/                   # 前端页面
│   ├── static/           # 静态资源
│   │   ├── css/          # 样式文件
│   │   ├── js/           # 脚本文件
│   │   └── favicon.ico   # 网站图标
│   ├── admin_index.html  # 管理后台首页
│   ├── admin_login.html  # 管理后台登录页
│   └── index.html        # 主页面
├── scripts/              # 脚本文件
├── log/                  # 日志文件
└── tool/                 # 工具脚本
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- 以下工具用于代码生成：
  - [wire](https://github.com/google/wire) - 依赖注入
  - [go-enum](https://github.com/abice/go-enum) - 枚举生成

### 安装依赖工具

```bash
# 安装 wire
go install github.com/google/wire/cmd/wire@latest

# 安装 go-enum
go install github.com/abice/go-enum@latest

# 注意：已移除 mockery，使用简化的测试方式
```

### 克隆项目

```bash
git clone <your-repo-url>
cd go-web-template
```

### 配置环境变量

创建 `.env` 文件：

```bash
# 应用配置
APP_SECRET=your-secret-key
APP_ADMIN_USER=admin
APP_ADMIN_PASS=admin123

# 数据库配置
DB_PASS=your-db-password

# JWT 配置
JWT_SECRET=your-jwt-secret
```

### 生成代码

```bash
# 生成依赖注入代码
make wire

# 生成枚举代码
make generate

# 注意：已移除 protobuf，使用标准 HTTP 接口
```

### 构建和运行

```bash
# 构建项目
make build

# 运行项目
make run

# 或者直接运行
go run cmd/main.go -c config/app.yaml -p 8081 -v
```

## 🔧 配置说明

### 应用配置 (config/app.yaml)

```yaml
verbose: true
env: "local"
app:
  name: "go-web-template"
  port: 8081
  web: "./web"
  secret: "${APP_SECRET}"
  admin_user: "${APP_ADMIN_USER}"
  admin_pass: "${APP_ADMIN_PASS}"
  debug: true
  log_level: "debug"
  log_file: "./log/app.log"
db:
  host: "localhost"
  port: 3306
  user: "root"
  pass: "${DB_PASS}"
  db: "mytest_db"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
jwt:
  secret: "${JWT_SECRET}"
  expire_time: "24h"
  issuer: "go-web-template"
```

## 📡 API 接口

### 基础 API

- `GET /` - 首页
- `GET /admin/` - 管理后台首页
- `GET /admin/login` - 管理后台登录页

### API 接口

- `POST /api/login` - 用户登录
- `GET /api/hello?name=World` - Hello 接口
- `GET /api/image/{key}` - 获取图片
- `POST /api/upload` - 上传图片

### 管理接口

- `POST /admin/login` - 管理员登录
- `POST /admin/add_token` - 生成访问令牌

## 🎛️ 管理后台

访问 `http://localhost:8081/admin/` 进入管理后台：

- **用户管理** - 生成和管理用户访问令牌
- **系统监控** - 查看系统状态和统计信息
- **日志查看** - 查看系统日志

默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/service/... -v

# 运行测试并生成覆盖率报告
go test ./... -cover
```

## 🐳 Docker 部署

```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run
```

## 📝 开发指南

### 添加新的 API 接口

1. 在 `internal/controller/request/` 中定义请求结构体
2. 在 `internal/controller/response/` 中定义响应结构体
3. 在控制器中实现 HTTP 处理函数
4. 在 `RegisterRoutes` 方法中注册路由
5. 使用 JWT 鉴权中间件

### JWT 认证系统

项目使用基于 JWT 的认证系统：

- **Bearer Token**: 用于 API 认证，通过 `Authorization: Bearer <token>` 头传递
- **Cookie 认证**: 用于 Web 界面认证，JWT token 存储在 cookie 中
- **配置驱动**: JWT 密钥和过期时间从 `config/app.yaml` 读取
- **角色管理**: 支持 `user` 和 `admin` 角色

### SQLite 数据库

项目使用 SQLite 数据库进行数据持久化：

- **自动初始化**: 首次启动时自动创建数据库表和默认管理员账户
- **默认账户**: 管理员账户 `admin/admin`，首次启动后请及时修改密码
- **用户管理**: 支持用户注册、登录、密码修改等功能
- **数据安全**: 密码使用哈希存储，确保数据安全

#### 数据库配置

```yaml
db:
  path: "./data/app.db"  # SQLite 数据库文件路径
```

#### 默认管理员账户

- **用户名**: admin
- **密码**: admin
- **角色**: admin

⚠️ **安全提醒**: 首次启动后请立即修改默认管理员密码！

### 添加新的服务

1. 在 `internal/service/` 中定义服务接口
2. 实现服务逻辑
3. 在 `internal/setup/` 中注册服务
4. 运行 `make wire` 生成依赖注入代码

### 添加新的数据模型

1. 在 `internal/model/` 中定义模型
2. 在 `internal/repository/` 中实现数据访问
3. 在 `internal/service/` 中使用仓库

## 🔍 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 查看端口占用
   lsof -i :8081
   # 杀死进程
   kill -9 <PID>
   ```

2. **依赖注入错误**
   ```bash
   # 重新生成 wire 代码
   make wire
   ```

3. **前端资源问题**
   ```bash
   # 检查静态资源路径
   ls -la web/static/
   ```

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请提交 Issue 或联系维护者。