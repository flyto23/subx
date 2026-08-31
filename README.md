# SublinkX - 订阅链接管理系统

一个轻量级的代理订阅链接管理系统，支持节点管理、模板配置和订阅生成。

## 核心功能

- ✅ **节点管理**: 支持多种协议 (VMess, VLESS, Trojan, Shadowsocks, Hysteria2 等)
- ✅ **模板管理**: Clash/Surge 配置文件模板
- ✅ **订阅管理**: 生成和管理订阅链接
- ✅ **客户端适配**: 自动识别客户端类型返回对应格式

## 快速开始

### 安装
```bash
curl -fsSL https://raw.githubusercontent.com/your-repo/sublinkx/main/install.sh | bash
```

### 运行
```bash
sublink run
```

### 管理菜单
```bash
sublink menu
```

## 目录结构

```
sublinkx/
├── api/           # API 接口层
├── models/        # 数据模型
├── node/          # 协议解析
├── routers/       # 路由配置
├── template/      # 配置模板
├── utils/         # 工具函数
├── webs/          # 前端源码
├── static/        # 静态资源
├── main.go        # 入口文件
├── install.sh     # 安装脚本
└── menu.sh        # 管理菜单
```

## 配置

配置文件位于 `./db/config.yaml`，支持环境变量覆盖。

## License

MIT
