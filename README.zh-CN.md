# RSS 监控器

一个用 Go 编写的轻量级 RSS 监控服务，运行在 GitHub Actions 上，通过 Bark 向 iOS 设备发送推送通知。

[English](README.md) | 简体中文

## 功能特性

- 🚀 自动运行在 GitHub Actions 上（免费额度）
- 📱 通过 Bark 向 iOS 设备推送通知
- 🔄 并发监控多个 RSS 源
- 💾 状态持久化在仓库中
- 🔍 智能去重（支持 GUID、链接或标题）
- 📦 单一静态二进制文件，无运行时依赖
- 🔧 易于配置和维护
- 🎯 支持单条和聚合通知
- 🌏 完整支持中文内容
- 🤖 AI 智能总结（可选，支持 OpenAI、Azure OpenAI 等）

## 快速开始

### 1. Fork 或克隆此仓库

```bash
git clone https://github.com/rsswatcher/rsswatcher.git
cd rsswatcher
```

### 2. 配置 RSS 源

复制示例配置并编辑：

```bash
cp feeds.yaml.example feeds.yaml
```

编辑 `feeds.yaml` 添加你的 RSS 源：

```yaml
feeds:
  - id: my-blog
    name: 我的博客
    url: https://example.com/rss.xml
    notify: true
    dedupe_key: guid
    aggregate: false
```

### 3. 设置 Bark

1. 在 iOS 设备上安装 [Bark](https://apps.apple.com/cn/app/bark-customed-notifications/id1403753865)
2. 打开应用并复制你的设备密钥（格式：`https://api.day.app/你的设备密钥`）
3. 设备密钥就是 URL 中 `/` 后面的部分

### 4. 配置 GitHub Secrets

进入你的仓库 Settings → Secrets and variables → Actions，添加：

- `BARK_DEVICE_KEY`：你的 Bark 设备密钥（必需）
- `BARK_SERVER`：自定义 Bark 服务器 URL（可选，默认为 `https://api.day.app`）

### 5. 启用 GitHub Actions

工作流配置为每 30 分钟运行一次。你也可以手动触发：

1. 进入仓库的 "Actions" 标签页
2. 选择 "RSS Monitor (Go + Bark)"
3. 点击 "Run workflow"

### 6. 提交并推送

```bash
git add feeds.yaml
git commit -m "添加我的 RSS 源"
git push
```

工作流将自动开始运行！

## 配置说明

### 源配置选项

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 源的唯一标识符 |
| `name` | string | 是 | 通知中显示的名称 |
| `url` | string | 是 | RSS/Atom 源的 URL |
| `notify` | boolean | 否 | 是否启用通知（默认：true） |
| `dedupe_key` | string | 否 | 去重键：`guid`、`link` 或 `title`（默认：`guid`） |
| `aggregate` | boolean | 否 | 是否发送聚合通知（默认：false） |
| `aggregate_window_minutes` | int | 否 | 聚合窗口时间（分钟，默认：30） |

### 配置示例

#### 单条通知

```yaml
feeds:
  - id: tech-blog
    name: 技术博客
    url: https://techblog.example.com/rss
    notify: true
    dedupe_key: guid
    aggregate: false
```

#### 聚合通知

```yaml
feeds:
  - id: news-feed
    name: 新闻源
    url: https://news.example.com/feed
    notify: true
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
```

## 本地开发

### 前置要求

- Go 1.22 或更高版本
- Git

### 构建

```bash
go build -o rsswatcher ./cmd/rsswatcher
```

### 本地运行

```bash
export BARK_DEVICE_KEY="你的设备密钥"
./rsswatcher --config feeds.yaml --state state/last_states.json
```

### 运行测试

```bash
go test ./...
```

## 架构

```
┌─────────────────┐
│ GitHub Actions  │
│   (定时调度)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   rsswatcher    │
│   (Go 二进制)    │
└────────┬────────┘
         │
    ┌────┴────────────────┐
    │                     │
    ▼                     ▼
┌─────────┐         ┌──────────┐
│  Feeds  │         │  State   │
│  .yaml  │         │  .json   │
└────┬────┘         └────┬─────┘
     │                   │
     │    ┌──────────────┘
     │    │
     ▼    ▼
┌──────────────┐
│   Fetcher    │
│   Parser     │
│   Deduper    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│    Bark      │
│  Notifier    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ iOS 设备     │
└──────────────┘
```

## 项目结构

```
.
├── cmd/
│   └── rsswatcher/
│       └── main.go           # 应用程序入口
├── internal/
│   ├── config/
│   │   └── config.go         # 配置解析
│   ├── fetcher/
│   │   └── fetcher.go        # RSS 源获取（带重试）
│   ├── parser/
│   │   └── parser.go         # RSS/Atom 解析
│   ├── deduper/
│   │   └── deduper.go        # 去重逻辑
│   ├── notifier/
│   │   └── bark.go           # Bark 通知客户端
│   └── state/
│       └── state.go          # 状态持久化
├── state/
│   └── last_states.json      # 状态文件（自动生成）
├── .github/
│   └── workflows/
│       └── rss-monitor.yml   # GitHub Actions 工作流
├── feeds.yaml                # 你的源配置
├── feeds.yaml.example        # 示例配置
├── go.mod
├── go.sum
└── README.md
```

## 故障排查

### 未收到通知

1. 检查 `BARK_DEVICE_KEY` 是否在 GitHub Secrets 中正确设置
2. 通过发送测试通知验证设备密钥：
   ```bash
   curl "https://api.day.app/你的设备密钥/测试/你好"
   ```
3. 检查 GitHub Actions 日志是否有错误

### 状态未更新

1. 确保工作流有仓库的写入权限
2. 检查分支是否受保护（可能需要额外配置）
3. 查看 Actions 日志中的 "Commit state" 步骤

### 源无法工作

1. 验证 RSS 源 URL 是否可访问
2. 检查源格式（RSS/Atom/JSON Feed）
3. 查看日志中的具体错误信息

### 工作流未运行

1. 检查仓库是否启用了 GitHub Actions
2. 验证 `.github/workflows/rss-monitor.yml` 中的 cron 调度
3. 注意：在仓库 60 天无活动后，定时工作流可能会被禁用

## AI 智能总结功能

RSS Watcher 支持使用大语言模型为文章生成智能总结。配置后，推送通知将包含 AI 生成的总结内容。

### 环境变量配置

**方式 1：使用 .env 文件（推荐用于本地开发）**

创建 `.env` 文件并配置：

```bash
cp .env.example .env
# 编辑 .env 文件，填入你的配置
```

```bash
API_ENDPOINT=https://api.openai.com/v1/chat/completions
API_KEY=your-api-key
MODEL_NAME=gpt-3.5-turbo
```

**方式 2：环境变量（用于生产环境）**

```bash
export API_ENDPOINT="https://api.openai.com/v1/chat/completions"
export API_KEY="your-api-key"
export MODEL_NAME="gpt-3.5-turbo"
```

### 支持的 API 服务

- OpenAI API
- Azure OpenAI Service
- 其他兼容 OpenAI API 格式的服务（如 Ollama、vLLM 等）

### 功能特点

- ✅ 可选功能：未配置时不启用，完全向后兼容
- ✅ 自动回退：API 调用失败时自动使用原始描述
- ✅ 中文优化：针对中文总结进行了优化

详细使用说明请参考：[AI 总结功能文档](docs/AI_SUMMARY.md)

## 高级用法

### 自定义 Bark 服务器

如果你运行自己的 Bark 服务器：

```bash
# 在 GitHub Secrets 中设置
BARK_SERVER=https://your-bark-server.com
```

### 调整调度频率

编辑 `.github/workflows/rss-monitor.yml`：

```yaml
on:
  schedule:
    - cron: '*/15 * * * *'  # 每 15 分钟
    # - cron: '0 * * * *'   # 每小时
    # - cron: '0 */6 * * *' # 每 6 小时
```

## 中文支持

本项目完全支持中文内容：

- ✅ 正确处理中文字符截断（按字符而非字节）
- ✅ 支持中文 RSS 源
- ✅ 中文通知标题和内容
- ✅ 中文配置名称
- ✅ 完整的中文文档

### 字符串截断说明

项目使用 Unicode 字符（rune）而不是字节来截断字符串，这确保：
- 中文字符不会被截断到一半
- 表情符号正确显示
- 其他多字节字符正确处理

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

MIT License - 详见 LICENSE 文件

## 安全

- 永远不要将 `BARK_DEVICE_KEY` 提交到仓库
- 使用 GitHub Secrets 存储所有敏感数据
- 检查状态文件确保没有存储敏感信息

## 致谢

- [gofeed](https://github.com/mmcdole/gofeed) - RSS/Atom 解析器
- [Bark](https://github.com/Finb/Bark) - iOS 通知服务
- 受到各种 RSS 监控解决方案的启发

## 常见问题

**问：这需要花费多少钱？**  
答：免费！GitHub Actions 为免费账户提供每月 2,000 分钟。

**问：可以用于 Android 吗？**  
答：可以调整通知器以使用其他服务，如 Telegram、Discord 或 Pushover。

**问：可以监控多少个源？**  
答：仅受 GitHub Actions 执行时间限制（通常最多约 6 小时）。大多数用户可以监控 50-100 个源。

**问：可以在私有仓库上工作吗？**  
答：可以，GitHub Actions 在公共和私有仓库上都可以工作。

## 支持

如果遇到任何问题，请：

1. 查看[故障排查](#故障排查)部分
2. 查看 GitHub Actions 日志
3. 提交 issue 并附上问题详情
