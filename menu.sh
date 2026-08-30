#!/bin/bash

# 全局常量
readonly INSTALL_DIR="/usr/local/bin/sublink"
readonly SERVICE_FILE="/etc/systemd/system/sublink.service"
readonly GITHUB_REPO="gooaclok819/sublinkX"

# 颜色定义
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

# 打印信息
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 root 权限
check_root() {
    if [ "$(id -u)" != "0" ]; then
        log_error "该脚本必须以 root 身份运行"
        exit 1
    fi
}

# 获取最新版本
get_latest_version() {
    curl --silent "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
        grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# 获取二进制文件名
get_binary_name() {
    local machine_type
    machine_type=$(uname -m)
    case "$machine_type" in
        x86_64) echo "sublink_amd64" ;;
        aarch64) echo "sublink_arm64" ;;
        *) log_error "不支持的机器类型：$machine_type"; exit 1 ;;
    esac
}

# 下载并更新程序
do_update() {
    local latest_release
    latest_release=$(get_latest_version)
    
    if [ -z "$latest_release" ]; then
        log_error "无法获取最新版本信息"
        return 1
    fi
    
    log_info "最新版本：$latest_release"
    
    local file_name
    file_name=$(get_binary_name)
    
    # 下载到临时目录
    cd /tmp || exit 1
    if ! curl -LO "https://github.com/${GITHUB_REPO}/releases/download/${latest_release}/${file_name}"; then
        log_error "下载失败"
        return 1
    fi
    
    chmod +x "$file_name"
    
    # 停止服务
    systemctl stop sublink 2>/dev/null || true
    
    # 移动文件
    mv -f "$file_name" "${INSTALL_DIR}/sublink" || {
        log_error "移动文件失败"
        return 1
    }
    
    log_info "更新完成"
    return 0
}

# 显示服务状态
show_status() {
    local latest_release status version
    latest_release=$(get_latest_version)
    status=$(systemctl is-active sublink 2>/dev/null || echo "inactive")
    version=$("${INSTALL_DIR}/sublink" --version 2>/dev/null || echo "未知")
    
    echo "================================"
    echo "最新版本：${latest_release}"
    echo "当前版本：${version}"
    echo "服务状态：${status}"
    echo "================================"
}

# 启动服务
start_service() {
    systemctl daemon-reload
    systemctl start sublink
    log_info "服务已启动"
}

# 停止服务
stop_service() {
    systemctl stop sublink
    systemctl daemon-reload
    log_info "服务已停止"
}

# 卸载服务
uninstall_service() {
    # 停止并禁用服务
    systemctl stop sublink 2>/dev/null || true
    systemctl disable sublink 2>/dev/null || true
    
    # 删除服务文件
    rm -f "$SERVICE_FILE"
    rm -f "${INSTALL_DIR}/sublink"
    rm -f "/usr/bin/sublink"
    
    read -p "是否删除模板文件和数据库？(y/n): " isDelete
    if [ "$isDelete" = "y" ]; then
        rm -rf "${INSTALL_DIR}/db"
        rm -rf "${INSTALL_DIR}/template"
        rm -rf "${INSTALL_DIR}/logs"
    fi
    
    systemctl daemon-reload
    log_info "卸载完成"
}

# 修改端口
change_port() {
    local Port
    read -p "请输入新的端口号：" Port
    
    if [ ! -f "$SERVICE_FILE" ]; then
        log_error "服务文件不存在：$SERVICE_FILE"
        return 1
    fi
    
    # 替换端口参数
    if grep -q "\-\-port" "$SERVICE_FILE"; then
        sed -i "s/--port [0-9]\+/--port ${Port}/" "$SERVICE_FILE"
    else
        sed -i "/^ExecStart=/ s|$| --port ${Port}|" "$SERVICE_FILE"
    fi
    
    systemctl daemon-reload
    systemctl restart sublink
    log_info "端口已修改为：${Port}，服务已重启"
}

# 重置账号密码
reset_credentials() {
    local User Password
    read -p "请输入新的账号：" User
    read -p "请输入新的密码：" Password
    
    if [ -z "$User" ] || [ -z "$Password" ]; then
        log_error "账号和密码不能为空"
        return 1
    fi
    
    "${INSTALL_DIR}/sublink" setting --username "$User" --password "$Password"
    systemctl restart sublink
    log_info "账号密码已重置"
}

# 主菜单
show_menu() {
    show_status
    echo ""
    echo "1. 启动服务"
    echo "2. 停止服务"
    echo "3. 卸载安装"
    echo "4. 查看服务状态"
    echo "5. 查看运行目录"
    echo "6. 修改端口"
    echo "7. 更新程序"
    echo "8. 重置账号密码"
    echo "0. 退出"
    echo ""
}

# 主函数
main() {
    check_root
    
    while true; do
        show_menu
        read -p "请选择一个选项：" option
        
        case $option in
            1) start_service ;;
            2) stop_service ;;
            3) uninstall_service ;;
            4) systemctl status sublink ;;
            5) 
                echo "运行目录：${INSTALL_DIR}"
                echo "需要备份的目录：db, template (模板文件可选)"
                cd "$INSTALL_DIR" || exit 1
                ;;
            6) change_port ;;
            7) do_update ;;
            8) reset_credentials ;;
            0) log_info "退出"; exit 0 ;;
            *) log_warn "无效的选项，请重新选择" ;;
        esac
        
        echo ""
        read -p "按回车键继续..."
    done
}

main
