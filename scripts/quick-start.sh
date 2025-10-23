#!/bin/bash

# Go Web 模板项目快速启动脚本
# 作者: Go Web Template Team
# 版本: 1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 命令未找到，请先安装 $1"
        exit 1
    fi
}

# 检查 Go 版本
check_go_version() {
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 1.21+"
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
    REQUIRED_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        print_error "Go 版本过低，需要 1.21+，当前版本: $GO_VERSION"
        exit 1
    fi
    
    print_success "Go 版本检查通过: $GO_VERSION"
}

# 安装依赖工具
install_tools() {
    print_info "检查并安装必要的开发工具..."
    
    # 检查并安装 wire
    if ! command -v wire &> /dev/null; then
        print_info "安装 wire..."
        go install github.com/google/wire/cmd/wire@latest
    else
        print_success "wire 已安装"
    fi
    
    # 检查并安装 go-enum
    if ! command -v go-enum &> /dev/null; then
        print_info "安装 go-enum..."
        go install github.com/abice/go-enum@latest
    else
        print_success "go-enum 已安装"
    fi
    
    # 注意：已移除 mockery，使用简化的测试方式
}

# 创建环境配置文件
create_env_file() {
    if [ ! -f .env ]; then
        print_info "创建环境配置文件..."
        cat > .env << EOF
# 应用配置
APP_SECRET=your-secret-key-$(date +%s)
APP_ADMIN_USER=admin
APP_ADMIN_PASS=admin123

# 数据库配置
DB_PASS=your-db-password

# JWT 配置
JWT_SECRET=your-jwt-secret-$(date +%s)
EOF
        print_success "环境配置文件已创建: .env"
        print_warning "请编辑 .env 文件，设置正确的配置值"
    else
        print_success "环境配置文件已存在"
    fi
}

# 生成代码
generate_code() {
    print_info "生成代码..."
    
    # 生成依赖注入代码
    if command -v wire &> /dev/null; then
        print_info "生成 wire 代码..."
        make wire
    fi
    
    # 生成枚举代码
    if command -v go-enum &> /dev/null; then
        print_info "生成枚举代码..."
        make generate
    fi
    
    # 注意：已移除 protobuf，使用标准 HTTP 接口
    
    print_success "代码生成完成"
}

# 构建项目
build_project() {
    print_info "构建项目..."
    make build
    print_success "项目构建完成"
}

# 启动项目
start_project() {
    print_info "启动项目..."
    print_info "服务将在 http://localhost:8081 启动"
    print_info "管理后台: http://localhost:8081/admin/"
    print_info "默认管理员账号: admin / admin123"
    print_info "按 Ctrl+C 停止服务"
    echo
    
    # 启动服务
    make run
}

# 显示帮助信息
show_help() {
    echo "Go Web 模板项目快速启动脚本"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -c, --check    仅检查环境，不启动服务"
    echo "  -b, --build    仅构建项目，不启动服务"
    echo "  -g, --generate 仅生成代码，不启动服务"
    echo "  -s, --start    直接启动服务（跳过环境检查）"
    echo
    echo "示例:"
    echo "  $0              # 完整流程：检查环境 -> 生成代码 -> 构建 -> 启动"
    echo "  $0 -c           # 仅检查环境"
    echo "  $0 -b           # 仅构建项目"
    echo "  $0 -s           # 直接启动服务"
}

# 主函数
main() {
    echo "🚀 Go Web 模板项目快速启动脚本"
    echo "=================================="
    echo
    
    # 解析命令行参数
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -c|--check)
            print_info "仅检查环境..."
            check_go_version
            install_tools
            print_success "环境检查完成"
            exit 0
            ;;
        -b|--build)
            print_info "仅构建项目..."
            check_go_version
            build_project
            exit 0
            ;;
        -g|--generate)
            print_info "仅生成代码..."
            check_go_version
            install_tools
            generate_code
            exit 0
            ;;
        -s|--start)
            print_info "直接启动服务..."
            start_project
            exit 0
            ;;
        "")
            # 默认流程
            ;;
        *)
            print_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
    
    # 完整流程
    print_info "开始完整启动流程..."
    
    # 1. 检查环境
    check_go_version
    install_tools
    
    # 2. 创建环境配置
    create_env_file
    
    # 3. 生成代码
    generate_code
    
    # 4. 构建项目
    build_project
    
    # 5. 启动服务
    start_project
}

# 捕获中断信号
trap 'print_warning "服务已停止"; exit 0' INT

# 运行主函数
main "$@"
