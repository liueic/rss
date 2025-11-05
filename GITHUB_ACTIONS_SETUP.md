# GitHub Actions 部署指南

本文档说明如何配置和运行 RSS Watcher 的 GitHub Actions 工作流。

## 📋 前置要求

1. 已 Fork 或克隆此仓库到你的 GitHub 账户
2. 已在 iOS 设备上安装 Bark 应用并获取设备密钥
3. 已配置 `feeds.yaml` 文件（包含你想要监控的 RSS 源）

## 🔐 步骤 1: 配置 GitHub Secrets

### 必需配置

1. 进入你的 GitHub 仓库
2. 点击 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret** 添加以下密钥：

#### BARK_DEVICE_KEY（必需）

- **名称**: `BARK_DEVICE_KEY`
- **值**: 从 Bark 应用复制的设备密钥
- **获取方式**:
  1. 在 iOS 设备上打开 Bark 应用
  2. 复制显示的 URL，格式类似：`https://api.day.app/YOUR_DEVICE_KEY/`
  3. `YOUR_DEVICE_KEY` 部分就是你要设置的密钥

### 可选配置

#### BARK_SERVER（可选）

- **名称**: `BARK_SERVER`
- **值**: 自定义 Bark 服务器 URL（如果使用自建服务器）
- **默认值**: 如果不设置，将使用 `https://api.day.app`

#### AI 总结功能（可选）

如果你启用了 AI 总结功能，需要配置：

- **API_ENDPOINT**: AI API 端点（例如：`https://api.openai.com/v1/chat/completions`）
- **API_KEY**: API 密钥
- **MODEL_NAME**: 模型名称（例如：`gpt-3.5-turbo`）

## ⚙️ 步骤 2: 配置工作流权限

确保工作流有权限提交代码到 `rss-state` 分支：

1. 进入仓库 **Settings** → **Actions** → **General**
2. 在 **Workflow permissions** 部分：
   - 选择 **Read and write permissions**
   - 勾选 **Allow GitHub Actions to create and approve pull requests**

## 🚀 步骤 3: 提交代码并触发构建

### 方式 1: 通过 Git 提交（推荐）

```bash
# 确保所有文件已提交
git add .
git commit -m "chore: setup GitHub Actions workflows"
git push origin main
```

推送后，GitHub Actions 会自动触发 `build.yml` 工作流进行构建。

### 方式 2: 手动触发构建

1. 进入仓库的 **Actions** 标签页
2. 在左侧选择 **Build RSS Watcher** 工作流
3. 点击右侧的 **Run workflow** 按钮
4. 选择分支（通常是 `main`）
5. 点击绿色的 **Run workflow** 按钮

## ✅ 步骤 4: 验证构建成功

1. 在 **Actions** 标签页查看 **Build RSS Watcher** 工作流
2. 等待构建完成（通常需要 1-2 分钟）
3. 确认构建成功（绿色 ✓）
4. 点击构建运行，在 **Artifacts** 部分应该能看到 `rsswatcher-binary`

## 🧪 步骤 5: 测试监控工作流

### 手动触发监控工作流

1. 在 **Actions** 标签页选择 **RSS Monitor (Go + Bark)** 工作流
2. 点击 **Run workflow** 按钮
3. 选择分支并运行

### 验证运行结果

1. 查看工作流日志，确认：
   - ✅ 成功下载或构建了二进制文件
   - ✅ 成功读取了 `feeds.yaml` 配置
   - ✅ 成功从 `rss-state` 分支读取了状态文件
   - ✅ RSS 源被成功处理
   - ✅ 如果有新内容，应该收到 Bark 通知

2. 检查 `rss-state` 分支：
   - 工作流会自动创建此分支（如果不存在）
   - 状态文件会保存在 `rss-state` 分支的 `state/last_states.json`

## 📅 自动运行

配置完成后，监控工作流会：
- **每 30 分钟自动运行一次**（根据 `monitor.yml` 中的 cron 配置）
- 自动从最新的构建 artifact 下载二进制文件
- 自动更新 `rss-state` 分支中的状态文件

## 🔍 故障排查

### 问题 1: 构建工作流失败

**可能原因**:
- Go 版本不兼容
- 代码语法错误
- 依赖问题

**解决方法**:
- 检查工作流日志中的错误信息
- 确保 `go.mod` 和 `go.sum` 文件已正确提交
- 尝试本地构建：`go build -o rsswatcher ./cmd/rsswatcher`

### 问题 2: 监控工作流无法下载 artifact

**可能原因**:
- 构建工作流尚未运行或失败
- Artifact 已过期（超过 7 天）

**解决方法**:
- 先运行 **Build RSS Watcher** 工作流
- 工作流会自动回退到构建二进制（如果 artifact 不存在）

### 问题 3: 没有收到 Bark 通知

**检查清单**:
1. ✅ `BARK_DEVICE_KEY` 是否正确配置
2. ✅ Bark 应用是否正常运行
3. ✅ 测试通知：在终端运行：
   ```bash
   curl "https://api.day.app/YOUR_DEVICE_KEY/Test/Hello"
   ```
4. ✅ 检查 `feeds.yaml` 中对应的 feed 是否设置了 `notify: true`
5. ✅ 查看工作流日志，确认是否有错误信息
6. ✅ RSS 源是否有新内容（如果没有新内容，不会发送通知）

### 问题 4: State 文件未更新

**可能原因**:
- 工作流没有写入权限
- `rss-state` 分支创建失败

**解决方法**:
- 检查工作流权限设置（见步骤 2）
- 查看工作流日志中的 "Commit and push state changes" 步骤
- 手动检查 `rss-state` 分支是否存在

### 问题 5: 工作流运行时间过长

**优化建议**:
- 减少监控的 RSS 源数量
- 检查是否有 RSS 源响应缓慢
- 考虑增加 cron 间隔时间（例如：从 30 分钟改为 1 小时）

## 📊 监控工作流状态

### 添加状态徽章

在你的 README.md 中添加工作流状态徽章：

```markdown
![Build Status](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/build.yml/badge.svg)
![Monitor Status](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/monitor.yml/badge.svg)
```

### 查看运行历史

- 进入 **Actions** 标签页
- 选择对应的工作流查看历史运行记录
- 点击任意运行查看详细日志

## 🔄 更新配置

### 修改 RSS 源

1. 编辑 `feeds.yaml` 文件
2. 提交并推送到 `main` 分支
3. 下次监控工作流运行时会自动使用新配置

### 修改运行频率

编辑 `.github/workflows/monitor.yml` 中的 cron 表达式：

```yaml
schedule:
  - cron: '*/30 * * * *'  # 每 30 分钟
  # - cron: '0 * * * *'    # 每小时
  # - cron: '0 */6 * * *'  # 每 6 小时
```

## 📝 注意事项

1. **Artifact 保留期**: 构建的二进制文件会保留 7 天，之后会自动删除。如果超过 7 天未运行构建工作流，监控工作流会自动构建二进制。

2. **State 分支**: `rss-state` 分支由工作流自动管理，建议不要手动修改此分支。

3. **GitHub Actions 免费额度**: 
   - 公开仓库：无限制
   - 私有仓库：每月 2,000 分钟免费
   - 典型使用：每次运行约 1-3 分钟，每月约 150-300 分钟

4. **定时任务延迟**: GitHub Actions 的定时任务可能会有 5-10 分钟的延迟，这是正常现象。

## 🎉 完成！

配置完成后，你的 RSS 监控服务就会自动运行了！每 30 分钟检查一次配置的 RSS 源，有新内容时会通过 Bark 发送通知到你的 iOS 设备。

如有问题，请查看工作流日志或参考本文档的故障排查部分。

