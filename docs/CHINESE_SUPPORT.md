# 中文支持说明 / Chinese Support

本文档说明 RSS Watcher 如何支持中文及其他 Unicode 字符。

This document explains how RSS Watcher supports Chinese and other Unicode characters.

## 主要改进 / Key Improvements

### 1. 字符截断（Character Truncation）

**问题 / Problem:**
- 原先使用字节长度截断字符串，会导致中文字符被截断到一半
- Previously used byte-length truncation, which would cut Chinese characters in half

**解决方案 / Solution:**
- 使用 `[]rune` 进行字符级别的截断
- 确保多字节字符（中文、日文、表情符号等）不会被破坏
- Use `[]rune` for character-level truncation
- Ensures multi-byte characters (Chinese, Japanese, emojis, etc.) are not corrupted

**代码示例 / Code Example:**

```go
// 之前 / Before - 会截断中文字符
func truncateOld(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..." // ❌ 可能导致乱码
}

// 现在 / Now - 正确处理中文字符
func truncate(s string, maxLen int) string {
    s = strings.TrimSpace(s)
    runes := []rune(s)
    if len(runes) <= maxLen {
        return s
    }
    return string(runes[:maxLen]) + "..." // ✅ 正确截断
}
```

### 2. 测试覆盖 / Test Coverage

完整的中文字符处理测试：
Comprehensive Chinese character handling tests:

```go
// 中文文本测试
{
    name:   "Chinese text longer than max",
    input:  "这是一个很长的中文句子需要被截断处理",
    maxLen: 10,
    want:   "这是一个很长的中文句...",
}

// 中英混合文本测试
{
    name:   "Mixed English and Chinese",
    input:  "Hello 世界 this is a test 测试",
    maxLen: 15,
    want:   "Hello 世界 this i...",
}

// 表情符号测试
{
    name:   "Emoji support",
    input:  "Hello 👋 World 🌍",
    maxLen: 10,
    want:   "Hello 👋 Wo...",
}
```

### 3. 受影响的模块 / Affected Modules

#### internal/notifier/bark.go
- ✅ `truncate()` - 通知标题和正文截断
- ✅ `truncateBytes()` - 字节级别的安全截断（可选）

#### internal/parser/parser.go
- ✅ `cleanDescription()` - RSS 描述清理和截断
- ✅ `cleanDescriptionBytes()` - 字节级别的描述截断（可选）

## 使用示例 / Usage Examples

### 中文 RSS 源配置 / Chinese RSS Feed Configuration

```yaml
feeds:
  - id: chinese-tech-blog
    name: 中文科技博客
    url: https://example.com/cn/rss
    notify: true
    dedupe_key: guid
    aggregate: false

  - id: chinese-news
    name: 新闻频道
    url: https://news.example.com/feed
    notify: true
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
```

### 通知示例 / Notification Examples

**单条通知 / Individual Notification:**
```
标题 / Title: [中文科技博客] Go 1.23 发布：新特性详解
正文 / Body: 本文介绍了 Go 1.23 版本的主要新特性，包括性能改进...
```

**聚合通知 / Aggregated Notification:**
```
标题 / Title: [新闻频道] 3 new items
正文 / Body:
科技新闻：AI 技术突破
经济报道：市场分析报告
体育资讯：国际比赛结果
```

## 技术细节 / Technical Details

### Unicode 字符处理 / Unicode Character Handling

Go 语言中的字符串是 UTF-8 编码的字节序列：
Strings in Go are UTF-8 encoded byte sequences:

- **byte**: 1 个字节 (1 byte)
- **rune**: 1 个 Unicode 码点 (1 Unicode code point)
  - 英文字母: 1 byte
  - 中文字符: 通常 3 bytes
  - 表情符号: 4+ bytes

### 示例对比 / Comparison Example

```go
text := "你好世界"

// 字节长度 / Byte length
len(text) // = 12 (每个中文字符 3 字节)

// 字符长度 / Character length
len([]rune(text)) // = 4 (4 个字符)
```

### 安全截断策略 / Safe Truncation Strategy

1. **字符截断（推荐）/ Character Truncation (Recommended)**
   ```go
   runes := []rune(text)
   truncated := string(runes[:maxChars])
   ```
   - ✅ 保证字符完整性
   - ✅ 适合显示和通知
   - ❌ 可能超过字节限制

2. **字节截断（特殊场景）/ Byte Truncation (Special Cases)**
   ```go
   for i := maxBytes; i > 0; i-- {
       if utf8.ValidString(text[:i]) {
           return text[:i] + "..."
       }
   }
   ```
   - ✅ 严格控制字节大小
   - ✅ 不会产生乱码
   - ❌ 可能截断更多字符

## 性能考虑 / Performance Considerations

转换为 `[]rune` 会有轻微的性能开销，但：
Converting to `[]rune` has a slight performance overhead, but:

- ✅ 对于通知消息等短文本，影响可忽略不计
- ✅ 正确性比性能更重要
- ✅ 避免了乱码问题带来的用户体验损失

For notification messages and short texts:
- Negligible impact
- Correctness is more important than performance
- Avoids poor user experience from corrupted text

## 测试方法 / Testing

运行中文支持测试：
Run Chinese support tests:

```bash
# 运行所有测试
go test ./...

# 运行通知器测试（包含中文测试）
go test ./internal/notifier -v

# 运行解析器测试
go test ./internal/parser -v
```

## 常见问题 / FAQ

### Q: 是否支持繁体中文？
**A:** 是的，完全支持。`rune` 处理所有 Unicode 字符，包括简体、繁体中文。

### Q: Does it support Traditional Chinese?
**A:** Yes, fully supported. `rune` handles all Unicode characters, including both Simplified and Traditional Chinese.

### Q: 其他语言（日文、韩文、阿拉伯文等）呢？
**A:** 完全支持所有 Unicode 字符，包括但不限于：
- 日文（Japanese）：ひらがな、カタカナ、漢字
- 韩文（Korean）：한글
- 阿拉伯文（Arabic）：العربية
- 表情符号（Emoji）：😀🎉🌍

### Q: What about other languages?
**A:** Fully supports all Unicode characters, including but not limited to:
- Japanese: Hiragana, Katakana, Kanji
- Korean: Hangul
- Arabic: العربية
- Emojis: 😀🎉🌍

### Q: 通知长度限制是多少？
**A:** 
- 标题：50 个字符
- 正文：100 个字符（单条）
- 聚合通知：每个标题 60 个字符，最多显示 5 条

### Q: What are the notification length limits?
**A:**
- Title: 50 characters
- Body: 100 characters (individual)
- Aggregated: 60 characters per title, max 5 items shown

## 相关文档 / Related Documentation

- [README 中文版](../README.zh-CN.md)
- [快速开始中文版](../QUICKSTART.zh-CN.md)
- [Go Unicode 文档](https://go.dev/blog/strings)

## 贡献 / Contributing

如果你发现中文支持的问题或有改进建议：
If you find issues with Chinese support or have improvements:

1. 提交 Issue / Open an issue
2. 提供示例 RSS 源 / Provide example RSS feed
3. 描述预期行为 / Describe expected behavior
4. 提交 PR / Submit a pull request

---

**注意 / Note:** 本项目的中文支持已经过充分测试，可以安全地用于生产环境。
The Chinese support in this project has been thoroughly tested and is safe for production use.
