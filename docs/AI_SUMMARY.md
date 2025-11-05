# AI 总结功能使用指南

## 功能概述

RSS Watcher 现在支持使用大语言模型（LLM）为 RSS 文章生成智能总结。当你配置了 API 信息后，推送通知将包含 AI 生成的总结内容，而不是原始的 RSS 描述。

## 环境变量配置

要启用 AI 总结功能，需要设置以下环境变量：

### 必需的环境变量

1. **API_ENDPOINT** - 大模型 API 的端点地址
   - OpenAI: `https://api.openai.com/v1/chat/completions`
   - Azure OpenAI: `https://your-resource.openai.azure.com/openai/deployments/your-deployment/chat/completions?api-version=2024-02-15-preview`
   - 其他兼容 OpenAI API 格式的服务

2. **API_KEY** - API 密钥
   - OpenAI: 你的 OpenAI API Key
   - Azure OpenAI: 你的 Azure API Key
   - 其他服务: 对应的 API Key

3. **MODEL_NAME** - 要使用的模型名称
   - OpenAI: `gpt-3.5-turbo`, `gpt-4`, `gpt-4-turbo-preview` 等
   - Azure OpenAI: 你的部署名称
   - 其他服务: 对应的模型名称

### 配置方式

#### 方式 1：使用 .env 文件（推荐用于本地开发）

创建 `.env` 文件（项目根目录下）：

```bash
# 复制示例文件
cp .env.example .env

# 编辑 .env 文件，填入你的配置
```

`.env` 文件示例：

```bash
# AI 总结功能配置
API_ENDPOINT=https://api.openai.com/v1/chat/completions
API_KEY=sk-your-openai-api-key
MODEL_NAME=gpt-3.5-turbo

# Bark 通知配置
BARK_DEVICE_KEY=your-bark-device-key
BARK_SERVER=https://api.day.app
```

**注意**：`.env` 文件已添加到 `.gitignore`，不会提交到版本控制系统。

#### 方式 2：环境变量（用于生产环境）

```bash
export API_ENDPOINT="https://api.openai.com/v1/chat/completions"
export API_KEY="sk-your-openai-api-key"
export MODEL_NAME="gpt-3.5-turbo"
```

### 配置示例

#### 使用 OpenAI

**使用 .env 文件：**
```bash
API_ENDPOINT=https://api.openai.com/v1/chat/completions
API_KEY=sk-your-openai-api-key
MODEL_NAME=gpt-3.5-turbo
```

**或使用环境变量：**
```bash
export API_ENDPOINT="https://api.openai.com/v1/chat/completions"
export API_KEY="sk-your-openai-api-key"
export MODEL_NAME="gpt-3.5-turbo"
```

#### 使用 Azure OpenAI

**使用 .env 文件：**
```bash
API_ENDPOINT=https://your-resource.openai.azure.com/openai/deployments/gpt-35-turbo/chat/completions?api-version=2024-02-15-preview
API_KEY=your-azure-api-key
MODEL_NAME=gpt-35-turbo
```

#### 使用本地部署的模型（如 Ollama、vLLM 等）

**使用 .env 文件：**
```bash
API_ENDPOINT=http://localhost:11434/v1/chat/completions
API_KEY=ollama
MODEL_NAME=llama2
```

## 工作流程

1. **未配置环境变量**：如果未设置上述环境变量，程序将使用原始的 RSS 描述内容，功能完全向后兼容。

2. **已配置环境变量**：
   - 程序会为每个新文章调用 AI API 生成总结
   - 如果 API 调用成功，推送通知将使用 AI 总结
   - 如果 API 调用失败（网络错误、API 错误等），将自动回退到原始描述

## 功能特点

- ✅ **可选功能**：未配置时不启用，不影响现有功能
- ✅ **自动回退**：API 调用失败时自动使用原始描述
- ✅ **智能截断**：总结内容自动截断到合适长度（200字符）
- ✅ **中文优化**：提示词针对中文总结进行了优化
- ✅ **兼容多种 API**：支持任何兼容 OpenAI API 格式的服务

## 使用示例

### 本地测试

**使用 .env 文件（推荐）：**

```bash
# 1. 创建 .env 文件
cp .env.example .env

# 2. 编辑 .env 文件，填入你的配置
# API_ENDPOINT=https://api.openai.com/v1/chat/completions
# API_KEY=sk-your-key
# MODEL_NAME=gpt-3.5-turbo
# BARK_DEVICE_KEY=your-bark-key

# 3. 运行程序（会自动加载 .env 文件）
./rsswatcher --config feeds.yaml --state state/last_states.json
```

**或使用环境变量：**

```bash
# 设置环境变量
export API_ENDPOINT="https://api.openai.com/v1/chat/completions"
export API_KEY="sk-your-key"
export MODEL_NAME="gpt-3.5-turbo"
export BARK_DEVICE_KEY="your-bark-key"

# 运行程序
./rsswatcher --config feeds.yaml --state state/last_states.json
```

### GitHub Actions 配置

在 GitHub Secrets 中添加：

- `API_ENDPOINT`: 你的 API 端点
- `API_KEY`: 你的 API 密钥
- `MODEL_NAME`: 模型名称

然后在 `.github/workflows/rss-monitor.yml` 中，这些环境变量会自动被使用。

## API 请求格式

程序使用标准的 OpenAI Chat Completions API 格式：

```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "user",
      "content": "请为以下文章生成一个简洁的中文总结..."
    }
  ],
  "temperature": 0.7,
  "max_tokens": 500
}
```

## 注意事项

1. **API 费用**：使用 AI 总结会产生 API 调用费用，请根据你的使用量选择合适的模型和计费方案。

2. **请求限制**：某些 API 服务可能有速率限制，如果遇到大量新文章，可能会触发限流。

3. **超时设置**：API 请求的超时时间为 30 秒，如果模型响应较慢可能会超时。

4. **内容长度**：输入内容限制在 8000 字符以内，超出部分会被截断。

5. **总结长度**：生成的总结限制在 500 tokens，通知中显示时会被截断到 200 字符。

## 故障排查

### 总结功能未启用

- 检查是否设置了所有三个必需的环境变量
- 检查环境变量名称是否正确（区分大小写）

### API 调用失败

- 检查 API 端点 URL 是否正确
- 检查 API Key 是否有效
- 检查网络连接是否正常
- 查看程序日志中的具体错误信息

### 总结内容为空

- 检查 API 响应格式是否正确
- 检查模型是否支持中文
- 查看程序日志了解详细错误

## 支持的 API 服务

理论上支持任何兼容 OpenAI Chat Completions API 格式的服务，包括但不限于：

- OpenAI API
- Azure OpenAI Service
- Google Vertex AI（需要适配）
- Anthropic Claude（需要适配）
- 本地部署的 Ollama
- 本地部署的 vLLM
- 其他兼容 OpenAI 格式的服务

