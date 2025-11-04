# Go Web 模板项目

本仓库提供一个基于 Gin 的全栈 Web 模板，集成了配置管理、身份认证、静态页面、后台任务和数据库访问。项目开箱即用，适合作为业务系统或内部管理平台的脚手架。

## 主要特性

- Gin 驱动的 HTTP 服务，统一注册 API、用户页面和管理后台
- 结合 `.env` 与 `config/app.yaml` 的分层配置，并通过 `tenz-io/gokit/cmd` 自动加载
- 内置自实现的 JWT 管理器，支持 Cookie 与 Bearer 两种鉴权方式
- 使用 GORM + SQLite 自动建表与默认管理员初始化
- Wire 提供依赖注入；go-enum 生成枚举辅助代码
- Cron 定时任务框架，预置健康检查示例
- Makefile 与 `scripts/quick-start.sh` 整合常用构建、运行、测试流程
- 前端静态资源位于 `web/`，包含登录、用户中心和管理后台示例页面

## 目录总览

```
.
├── cmd/                   # 程序入口与 Server 装配
│   └── webgo/             # HTTP 服务器实现
├── config/                # YAML 配置（支持 env 占位符）
├── data/                  # SQLite 数据文件目录（首次运行自动创建）
├── docs/                  # 额外文档（如 API 结构说明）
├── internal/              # 业务代码
│   ├── config/            # 配置结构体
│   ├── controller/        # API、用户、管理后台控制器
│   ├── database/          # GORM 封装与初始化
│   ├── job/               # 定时任务
│   ├── middleware/        # 日志、鉴权等中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层（DAO）
│   ├── service/           # 业务逻辑
│   └── setup/             # Wire 依赖注入装配
├── scripts/               # 辅助脚本与 HTTP 调试文件
├── tool/                  # 可选工具构建目标
├── web/                   # 静态页面与前端资源
└── Makefile               # 构建、运行、测试命令集合
```

## 快速开始

### 方式一：使用快速启动脚本

```bash
make quick-start
```

脚本会检查 Go 环境、安装 `wire`/`go-enum`、生成代码、构建并运行服务。默认服务地址为 `http://localhost:8090`，日志写入 `./log/app.log`。

### 方式二：手动执行

1. **准备环境**
   - Go 1.21+（`go.mod` 目标版本为 1.24）
   - 代码生成工具：
     ```bash
     go install github.com/google/wire/cmd/wire@latest
     go install github.com/abice/go-enum@latest
     ```

2. **克隆仓库并进入目录**
   ```bash
   git clone <your-repo-url>
   cd go-web-template
   ```

3. **配置环境变量**
   项目默认读取 `.env`，至少需要设置 JWT 密钥：
   ```bash
   cat > .env <<'EOF'
   JWT_SECRET=please-change-me
   EOF
   ```

4. **生成依赖代码**
   ```bash
   make wire
   make generate
   ```

5. **构建并运行**
   ```bash
   make build
   make run    # 或 go run cmd/main.go -c config/app.yaml -p 8090 -v
   ```

## 配置说明

应用配置位于 `config/app.yaml`，支持 `${ENV_NAME}` 占位符。关键字段如下：

```yaml
verbose: true          # 控制 Gin 模式与日志输出
env: "local"
app:
  name: "go-web-template"
  port: 8090           # 启动端口（Makefile 运行时会覆盖为 8090）
  web: "./web"         # 静态资源目录
  debug: true          # 影响内部调试开关
  log_level: "debug"
  log_file: "./log/app.log"
db:
  path: "./data/app.db"
  debug: false
jwt:
  secret: "${JWT_SECRET}"
  expire_time: "24h"
  issuer: "go-web-template"
```

> 注意：`db` 配置中的连接池参数目前未使用，如需扩展可在 `internal/database` 中补充。

## 数据库与默认账户

- 首次启动时会自动在 `db.path` 指定位置创建 SQLite 数据库，并执行数据表迁移。
- 若不存在管理员账号，会初始化 `admin/admin`（固定盐值，仅用于开发演示）。请在生产环境启动后立即修改密码。
- 密码以 HMAC-SHA256 + Salt 的方式存储。

## Web 界面与 API

### 页面

- `GET /`：首页示例
- `GET /login`：登录页
- `GET /user/home`：用户中心（需要登录）
- `GET /admin/home`：管理后台（需要管理员登录）

### 认证流程

- `POST /login`：使用 JSON (`username`/`password`) 登录，成功后下发 Cookie `jwt_token`
- `POST /logout`：退出登录，清除 Cookie
- `POST /auth/change_password`：修改密码，需 Cookie 鉴权

### API 端点

| 模块   | 方法  | 路径                    | 说明                         | 鉴权方式                 |
| ------ | ----- | ----------------------- | ---------------------------- | ------------------------ |
| 用户   | POST  | `/user/generate_token`  | 生成自定义过期时间的 API Token | Cookie，角色 ≥ user      |
| 管理   | GET   | `/admin/users`          | 查询用户列表                 | Cookie，角色 = admin     |
| 管理   | POST  | `/admin/add_user`       | 新建用户                     | Cookie，角色 = admin     |
| 管理   | DELETE | `/admin/delete_user`    | 删除用户                     | Cookie，角色 = admin     |
| API    | GET   | `/api/hello?name=`      | 示例接口，返回数据库中的用户信息 | Bearer JWT，角色 ≥ user |

Bearer Token 可以通过 `/user/generate_token` 获取，Cookie 由 `/login` 下发。

`scripts/http/*.http` 提供了 VS Code REST Client 风格的调用示例，便于本地调试。

## 后台任务

`internal/job` 使用 `robfig/cron` 注册定时任务，示例 `HealthReporter` 每周六 23:40 触发健康报告生成。可在此基础上扩展业务任务。（如需停止或关闭，调用 `cron.Stop()` 并清理资源。）

## 开发指南

- **依赖注入**：在新增 service/repository/controller 后，记得将构造函数写入 `internal/setup/provider.go` 并执行 `make wire` 重新生成 `wire_gen.go`。
- **枚举生成**：在枚举定义上添加 `//go:generate go-enum --marshal`，运行 `make generate` 即可生成 `<name>_enum.go`。
- **控制器结构**：请求与响应结构定义在 `internal/controller/request` 与 `internal/controller/response`，详情见 `docs/API_STRUCTURE.md`。
- **静态资源**：`web/` 内包含 HTML、CSS 与 JS 示例，Gin 会根据 `app.web` 路径加载模板与静态资源。
- **工具链**：`tool/` 目录可存放自定义命令，使用 `make build-tools` 构建。

## 测试

```bash
make test               # 运行全部测试并输出覆盖率
go test ./... -cover    # 自定义命令
```

## 常用 Make 命令

- `make help`：查看所有可用命令
- `make init`：安装工具并生成代码
- `make dev`：生成代码、构建并以开发模式运行
- `make run`：构建后运行服务
- `make clean`：清理构建产物与日志

## 许可证与贡献

- 许可证：MIT，详情见 `LICENSE`
- 欢迎通过 Issue 或 Pull Request 反馈问题、贡献代码

如需定制化或遇到问题，可先查阅日志与配置，仍无法解决时请联系维护者。祝使用愉快！
