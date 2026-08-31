#!/bin/bash

# 构建脚本 - 支持多平台编译

set -e  # 遇到错误立即退出

# 颜色定义
GREEN='\033[0;32m'
NC='\033[0m'
log_info() { echo -e "${GREEN}[BUILD]${NC} $1"; }

# 清理旧的构建文件
log_info "清理旧的构建文件..."
rm -f sublink_amd64 sublink_arm64 sublink.exe

# 编译 Linux AMD64
log_info "编译 Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o sublink_amd64 main.go

# 编译 Linux ARM64
log_info "编译 Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o sublink_arm64 main.go

# 编译 Windows AMD64
log_info "编译 Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o sublink.exe main.go

log_info "构建完成!"
ls -lh sublink_*
