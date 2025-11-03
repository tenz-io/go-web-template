# Go Web 模板项目 Makefile
# 作者: Go Web Template Team
# 版本: 2.0.0

# 项目配置
REPO_NAME := go-web-template
BIN_DIR := bin
LOG_DIR := log
CONFIG_FILE := config/app.yaml
PORT := 8090

GOCACHE := $(CURDIR)/.cache/go-build
export GOCACHE

# 颜色定义
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
RED := \033[0;31m
NC := \033[0m # No Color

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "$(BLUE)Go Web 模板项目构建工具$(NC)"
	@echo "=========================="
	@echo ""
	@echo "$(BLUE)可用命令:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BLUE)示例:$(NC)"
	@echo "  make help          # 显示此帮助信息"
	@echo "  make quick-start   # 快速启动（推荐）"
	@echo "  make dev           # 开发模式启动"
	@echo "  make build         # 构建项目"
	@echo "  make test          # 运行测试"

# 快速启动（推荐）
.PHONY: quick-start
quick-start: ## 快速启动项目（推荐）
	@echo "$(GREEN)[INFO]$(NC) 启动快速启动脚本..."
	@chmod +x scripts/quick-start.sh
	@./scripts/quick-start.sh

# 开发模式
.PHONY: dev
dev: generate build ## 开发模式启动
	@echo "$(GREEN)[INFO]$(NC) 启动开发模式..."
	@echo "$(YELLOW)[INFO]$(NC) 服务地址: http://localhost:$(PORT)"
	@echo "$(YELLOW)[INFO]$(NC) 管理后台: http://localhost:$(PORT)/admin/"
	@echo "$(YELLOW)[INFO]$(NC) 按 Ctrl+C 停止服务"
	@./$(BIN_DIR)/$(REPO_NAME) -c $(CONFIG_FILE) -p $(PORT) -v

# 环境检查
.PHONY: check
check: ## 检查开发环境
	@echo "$(BLUE)[INFO]$(NC) 检查 Go 版本..."
	@go version
	@echo "$(BLUE)[INFO]$(NC) 检查依赖工具..."
	@which wire || echo "$(YELLOW)[WARNING]$(NC) wire 未安装，运行: go install github.com/google/wire/cmd/wire@latest"
	@which go-enum || echo "$(YELLOW)[WARNING]$(NC) go-enum 未安装，运行: go install github.com/abice/go-enum@latest"
	@echo "$(GREEN)[SUCCESS]$(NC) 环境检查完成"

# 安装依赖
.PHONY: deps
deps: ## 安装项目依赖
	@echo "$(BLUE)[INFO]$(NC) 安装项目依赖..."
	@go mod tidy -v
	@echo "$(GREEN)[SUCCESS]$(NC) 依赖安装完成"

# 安装开发工具
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "$(BLUE)[INFO]$(NC) 安装开发工具..."
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/abice/go-enum@latest
	@echo "$(GREEN)[SUCCESS]$(NC) 开发工具安装完成"

# 生成代码
.PHONY: generate
generate: ensure-cache ## 生成所有代码
	@echo "$(BLUE)[INFO]$(NC) 生成枚举代码..."
	@go generate ./internal/constant/... ./internal/repository/... ./internal/service/...
	@echo "$(GREEN)[SUCCESS]$(NC) 代码生成完成"

# 生成 wire 代码
.PHONY: wire
wire: ensure-cache ## 生成依赖注入代码
	@echo "$(BLUE)[INFO]$(NC) 生成 wire 代码..."
	@wire gen $(REPO_NAME)/internal/setup/...
	@echo "$(GREEN)[SUCCESS]$(NC) Wire 代码生成完成"

# 代码格式化
.PHONY: fmt
fmt: ## 格式化代码
	@echo "$(BLUE)[INFO]$(NC) 格式化代码..."
	@go fmt ./...
	@echo "$(GREEN)[SUCCESS]$(NC) 代码格式化完成"

# 代码检查
.PHONY: lint
lint: ## 运行代码检查
	@echo "$(BLUE)[INFO]$(NC) 运行代码检查..."
	@go vet ./...
	@echo "$(GREEN)[SUCCESS]$(NC) 代码检查完成"

# 构建项目
.PHONY: build
build: wire ## 构建项目
	@echo "$(BLUE)[INFO]$(NC) 构建项目..."
	@mkdir -p $(BIN_DIR)
	@go build -mod=readonly -v -o $(BIN_DIR)/$(REPO_NAME) ./cmd
	@echo "$(GREEN)[SUCCESS]$(NC) 项目构建完成: $(BIN_DIR)/$(REPO_NAME)"

# 运行项目
.PHONY: run
run: build ## 构建并运行项目
	@echo "$(BLUE)[INFO]$(NC) 启动服务..."
	@echo "$(YELLOW)[INFO]$(NC) 服务地址: http://localhost:$(PORT)"
	@echo "$(YELLOW)[INFO]$(NC) 管理后台: http://localhost:$(PORT)/admin/"
	@echo "$(YELLOW)[INFO]$(NC) 按 Ctrl+C 停止服务"
	@./$(BIN_DIR)/$(REPO_NAME) -c $(CONFIG_FILE) -p $(PORT) -v

# 测试
.PHONY: test
test: ## 运行测试
	@echo "$(BLUE)[INFO]$(NC) 运行测试..."
	@go test -v ./... -cover
	@echo "$(GREEN)[SUCCESS]$(NC) 测试完成"

# 测试覆盖率
.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "$(BLUE)[INFO]$(NC) 运行测试覆盖率分析..."
	@go test -v ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)[SUCCESS]$(NC) 覆盖率报告已生成: coverage.html"

# 清理
.PHONY: clean
clean: ## 清理构建文件
	@echo "$(BLUE)[INFO]$(NC) 清理构建文件..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(LOG_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)[SUCCESS]$(NC) 清理完成"

# Docker 构建
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "$(BLUE)[INFO]$(NC) 构建 Docker 镜像..."
	@docker build --build-arg SKAFFOLD_GO_GCFLAGS="-N -l" -t $(REPO_NAME) .
	@echo "$(GREEN)[SUCCESS]$(NC) Docker 镜像构建完成: $(REPO_NAME)"

# Docker 运行
.PHONY: docker-run
docker-run: docker-build ## 构建并运行 Docker 容器
	@echo "$(BLUE)[INFO]$(NC) 启动 Docker 容器..."
	@docker rm -f $(REPO_NAME) || true
	@docker run --name $(REPO_NAME) -p $(PORT):$(PORT) $(REPO_NAME)
	@echo "$(GREEN)[SUCCESS]$(NC) Docker 容器启动完成"

# 创建环境配置
.PHONY: env
env: ## 创建环境配置文件
	@echo "$(BLUE)[INFO]$(NC) 创建环境配置文件..."
	@if [ ! -f .env ]; then \
		echo "# 应用配置" > .env; \
		echo "APP_SECRET=your-secret-key-$$(date +%s)" >> .env; \
		echo "APP_ADMIN_USER=admin" >> .env; \
		echo "APP_ADMIN_PASS=admin123" >> .env; \
		echo "" >> .env; \
		echo "# 数据库配置" >> .env; \
		echo "DB_PASS=your-db-password" >> .env; \
		echo "" >> .env; \
		echo "# JWT 配置" >> .env; \
		echo "JWT_SECRET=your-jwt-secret-$$(date +%s)" >> .env; \
		echo "$(GREEN)[SUCCESS]$(NC) 环境配置文件已创建: .env"; \
	else \
		echo "$(YELLOW)[INFO]$(NC) 环境配置文件已存在"; \
	fi

# 显示项目信息
.PHONY: info
info: ## 显示项目信息
	@echo "$(BLUE)[INFO]$(NC) 项目信息:"
	@echo "  项目名称: $(REPO_NAME)"
	@echo "  配置文件: $(CONFIG_FILE)"
	@echo "  默认端口: $(PORT)"
	@echo "  二进制文件: $(BIN_DIR)/$(REPO_NAME)"
	@echo ""
	@echo "$(BLUE)[INFO]$(NC) 服务地址:"
	@echo "  主页: http://localhost:$(PORT)/"
	@echo "  管理后台: http://localhost:$(PORT)/admin/"
	@echo "  API 接口: http://localhost:$(PORT)/api/"
	@echo ""
	@echo "$(BLUE)[INFO]$(NC) 默认管理员账号:"
	@echo "  用户名: admin"
	@echo "  密码: admin123"

# 完整构建流程
.PHONY: all
all: clean deps generate build test ## 完整构建流程

.PHONY: ensure-cache
ensure-cache:
	@mkdir -p $(GOCACHE)
	@echo "$(GREEN)[SUCCESS]$(NC) 完整构建流程完成"

# 开发环境初始化
.PHONY: init
init: install-tools deps generate wire env ## 初始化开发环境
	@echo "$(GREEN)[SUCCESS]$(NC) 开发环境初始化完成"

# 默认目标
.DEFAULT_GOAL := help
