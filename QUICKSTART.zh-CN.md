# 快速开始指南

5 分钟内完成 RSS Watcher 的部署！

[English](QUICKSTART.md) | 简体中文

## 1. 获取 Bark 设备密钥（1 分钟）

1. 在 iPhone 上安装 [Bark](https://apps.apple.com/cn/app/bark-customed-notifications/id1403753865)
2. 打开应用
3. 复制你的设备密钥（URL 中显示的长字符串）

## 2. Fork 此仓库（30 秒）

点击本页面右上角的 "Fork" 按钮。

## 3. 添加你的 Bark 密钥（30 秒）

在你 fork 的仓库中：

1. 进入 **Settings（设置）** → **Secrets and variables** → **Actions**
2. 点击 **New repository secret（新建仓库密钥）**
3. 名称：`BARK_DEVICE_KEY`
4. 值：[粘贴你的设备密钥]
5. 点击 **Add secret（添加密钥）**

## 4. 添加你的 RSS 源（2 分钟）

1. 在仓库中点击 `feeds.yaml` 文件
2. 点击铅笔图标（编辑）
3. 替换内容为：

```yaml
feeds:
  - id: my-first-feed
    name: 我的博客
    url: https://example.com/rss.xml  # 替换为你的 RSS 源 URL
    notify: true
    dedupe_key: guid
    aggregate: false
```

4. 点击 **Commit changes（提交更改）**

## 5. 测试运行！（1 分钟）

1. 进入 **Actions（操作）** 标签页
2. 点击左侧的 **RSS Monitor (Go + Bark)**
3. 点击右上角的 **Run workflow（运行工作流）** 按钮
4. 点击绿色的 **Run workflow（运行工作流）** 按钮
5. 等待约 30 秒
6. 查看你的 iPhone 是否收到通知！🎉

## 完成！

你的 RSS 监控器现在将每 30 分钟自动运行一次。

## 下一步

- 在 `feeds.yaml` 中添加更多源
- 阅读完整的 [README.zh-CN.md](README.zh-CN.md) 了解所有功能
- 查看 [DEPLOYMENT.md](DEPLOYMENT.md) 了解高级配置（英文）

## 故障排查

**没有收到通知？**

1. 测试 Bark：`curl "https://api.day.app/你的密钥/测试/你好"`
2. 检查 Actions 日志是否有错误
3. 验证你的源 URL 在浏览器中是否可以访问

**需要帮助？** [提交 issue](https://github.com/rsswatcher/rsswatcher/issues)

## 中文内容支持

本项目完全支持中文 RSS 源和通知：

✅ 中文标题和描述正确显示  
✅ 中文字符不会被截断到一半  
✅ 支持中文配置  
✅ 完整的中文文档  

示例中文配置：

```yaml
feeds:
  - id: tech-news-cn
    name: 科技新闻
    url: https://example.com/cn/rss
    notify: true
    dedupe_key: guid
    aggregate: false
    
  - id: blog-cn
    name: 技术博客
    url: https://blog.example.com/feed
    notify: true
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
```
