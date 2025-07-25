#  MiniRustpbxgo

一个集成 WebRTC、Go 后端和 Rust 服务的实时语音通信系统，支持无缝音频传输、语音处理和 AI 交互功能。

## 概述

MiniRustpbxgo 结合了 WebRTC 技术与基于 Go 的信令服务器和 Rust 后端服务，打造了一个强大的语音通信平台。它支持实时音频流传输、语音识别 (ASR)、文本转语音 (TTS)，以及通过 LLM 集成实现的 AI 驱动响应。

## 核心功能

- **WebRTC 集成**：实时点对点音频通信，包含 ICE 候选者管理
- **WebSocket 信令**：用于会话建立和消息路由的安全 WebSocket 连接
- **语音识别 (ASR)**：实时和最终的语音转文本转换
- **文本转语音 (TTS)**：高质量音频生成，支持可配置参数
- **AI 交互**：与 LLM 模型集成，提供智能响应
- **多提供商支持**：兼容多种 ASR/TTS 提供商（腾讯、阿里云等）
- **机器人管理**：可配置机器人配置文件，支持自定义声音和行为

## 技术栈

- **前端**：HTML5、JavaScript（WebRTC API）
- **后端**：Go 1.24+ 搭配 Gin 框架
- **服务**：Rust 集成用于性能关键型操作
- **数据库**：MySQL 用于持久化存储
- **缓存**：Redis 用于会话管理
- **网络**：WebSocket 用于信令，WebRTC 用于媒体流传输
- **AI**：OpenAI API 集成用于 LLM 能力

## 前置要求

- Go 1.24 或更高版本
- Rust 工具链
- MySQL 数据库
- Redis 服务器
- 支持 WebRTC 的现代 Web 浏览器
- 有效的 ASR/TTS 提供商 API 密钥（如使用云服务）

## 安装步骤

- 克隆仓库：

```bash
git clone https://github.com/yourusername/miniRustpbxgo.git
cd miniRustpbxgo
```

- 安装 Go 依赖：

```bash
go mod download
```

- 配置环境变量：

    - ·数据库连接字符串

    - Redis 服务器地址

    - ASR/TTS 提供商 API 密钥

    - LLM API 凭证

- 构建并运行应用：

```bash
go build -o miniRustpbxgo cmd/main.go
./miniRustpbxgo
```

- 在 Web 浏览器中打开index.html访问客户端界面

## 使用方法

1. 点击 "建立语音连接" 初始化 WebRTC 会话
2. 对着麦克风说话进行语音识别
3. 接收实时转录文本和 AI 生成的响应
4. 通过机器人配置调整 TTS 参数（语速、音量、发音人）
5. 点击 "挂断通话" 终止连接

## API 端点

- `GET /health`：服务健康检查
- `POST /robot/create`：创建新的机器人配置
- `WS /out/webrtc/setup`：用于 WebRTC 信令的 WebSocket 端点

## 配置说明

机器人配置和 API 密钥存储在数据库中。使用`RobotCreate` API 配置新的机器人配置文件，可自定义：

- 语速（0.5-2.0）
- 音量级别（0-10）
- 发音人声音
- 情感语调
- AI 行为的系统提示词

## 贡献指南

1. Fork 本仓库
2. 创建特性分支（`git checkout -b feature/amazing-feature`）
3. 提交更改（`git commit -m 'Add some amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 打开 Pull Request

## 故障排除

- 确保 WebSocket 连接未被防火墙阻止
- 验证生产环境部署的 SSL 配置
- 检查 ASR/TTS 服务的 API 密钥有效性
- 确认 WebRTC 功能的正确 CORS 设置