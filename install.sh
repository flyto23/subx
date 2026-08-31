#!/bin/bash

# 全局常量
readonly INSTALL_DIR="/usr/local/bin/sublink"
readonly GITHUB_REPO="gooaclok819/sublinkX"
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly NC='\033[0m'

# 打印信息
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 root 权限
if [ "$(id -u)" != "0" ]; then
    log_error "该脚本必须以 root 身份运行"
    exit 1
fi

# 创建安装目录
mkdir -p "$INSTALL_DIR"

# 获取最新版本和二进制文件名
latest_release=$(curl --silent "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
log_info "最新版本：$latest_release"

machine_type=$(uname -m)
case "$machine_type" in
    x86_64) file_name="sublink_amd64" ;;
    aarch64) file_name="sublink_arm64" ;;
    *) log_error "不支持的机器类型：$machine_type"; exit 1 ;;
esac

# 下载并安装
cd ~ || exit 1
curl -LO "https://github.com/${GITHUB_REPO}/releases/download/${latest_release}/${file_name}"
chmod +x "$file_name"
mv -f "$file_name" "${INSTALL_DIR}/sublink"

# 创建 systemd 服务
cat > /etc/systemd/system/sublink.service << EOFSystemService
[Unit]
Description=Sublink Service
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/sublink run
WorkingDirectory=${INSTALL_DIR}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOFSystemService

# 启动服务
systemctl daemon-reload
systemctl start sublink
systemctl enable sublink

log_info "服务已启动并已设置为开机启动"
log_info "默认账号：admin 密码：123456 端口：8000"
log_info "输入 sublink 可以呼出管理菜单"

# 安装管理脚本
curl -o /usr/bin/sublink -H "Cache-Control: no-cache" -H "Pragma: no-cache" \
    https://raw.githubusercontent.com/${GITHUB_REPO}/main/menu.sh
chmod 755 /usr/bin/sublink
