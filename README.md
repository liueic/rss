# RSS Watcher

A lightweight RSS monitoring service written in Go that runs on GitHub Actions and sends iOS notifications via Bark.

English | [简体中文](README.zh-CN.md)

## Features

- 🚀 Runs automatically on GitHub Actions (free tier)
- 📱 Push notifications to iOS devices via Bark
- 🔄 Monitors multiple RSS feeds concurrently
- 💾 State persistence in repository
- 🔍 Smart deduplication (by GUID, link, or title)
- 📦 Single static binary, no runtime dependencies
- 🔧 Easy to configure and maintain
- 🎯 Support for both individual and aggregated notifications
- 🌏 Full Unicode support (Chinese, Japanese, etc.)

## Quick Start

### 1. Fork or Clone This Repository

```bash
git clone https://github.com/rsswatcher/rsswatcher.git
cd rsswatcher
```

### 2. Configure Your Feeds

Copy the example configuration and edit it:

```bash
cp feeds.yaml.example feeds.yaml
```

Edit `feeds.yaml` to add your RSS feeds:

```yaml
feeds:
  - id: my-blog
    name: My Favorite Blog
    url: https://example.com/rss.xml
    notify: true
    dedupe_key: guid
    aggregate: false
```

### 3. Set Up Bark

1. Install [Bark](https://apps.apple.com/app/bark-customed-notifications/id1403753865) on your iOS device
2. Open the app and copy your device key (format: `https://api.day.app/YOUR_DEVICE_KEY`)
3. The device key is the part after `/` in the URL

### 4. Configure GitHub Secrets

Go to your repository Settings → Secrets and variables → Actions, and add:

- `BARK_DEVICE_KEY`: Your Bark device key (required)
- `BARK_SERVER`: Custom Bark server URL (optional, defaults to `https://api.day.app`)

### 5. Enable GitHub Actions

The workflow is configured to run every 30 minutes. You can also trigger it manually:

1. Go to the "Actions" tab in your repository
2. Select "RSS Monitor (Go + Bark)"
3. Click "Run workflow"

### 6. Commit and Push

```bash
git add feeds.yaml
git commit -m "Add my RSS feeds"
git push
```

The workflow will start running automatically!

## Configuration

### Feed Configuration Options

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the feed |
| `name` | string | Yes | Display name for notifications |
| `url` | string | Yes | RSS/Atom feed URL |
| `notify` | boolean | No | Enable notifications (default: true) |
| `dedupe_key` | string | No | Deduplication key: `guid`, `link`, or `title` (default: `guid`) |
| `aggregate` | boolean | No | Send aggregated notifications (default: false) |
| `aggregate_window_minutes` | int | No | Aggregation window in minutes (default: 30) |

### Example Configurations

#### Individual Notifications

```yaml
feeds:
  - id: tech-blog
    name: Tech Blog
    url: https://techblog.example.com/rss
    notify: true
    dedupe_key: guid
    aggregate: false
```
#### Aggregated Notifications

```yaml
feeds:
  - id: news-feed
    name: News Feed
    url: https://news.example.com/feed
    notify: true
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
```

## Local Development

### Prerequisites

- Go 1.22 or later
- Git

### Build

```bash
go build -o rsswatcher ./cmd/rsswatcher
```

### Run Locally

```bash
export BARK_DEVICE_KEY="your-device-key"
./rsswatcher --config feeds.yaml --state state/last_states.json
```

### Run Tests

```bash
go test ./...
```

## Architecture

```
┌─────────────────┐
│ GitHub Actions  │
│   (Scheduler)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   rsswatcher    │
│   (Go Binary)   │
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
│ iOS Device   │
└──────────────┘
```

## Project Structure

```
.
├── cmd/
│   └── rsswatcher/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration parsing
│   ├── fetcher/
│   │   └── fetcher.go        # RSS feed fetching with retry
│   ├── parser/
│   │   └── parser.go         # RSS/Atom parsing
│   ├── deduper/
│   │   └── deduper.go        # Deduplication logic
│   ├── notifier/
│   │   └── bark.go           # Bark notification client
│   └── state/
│       └── state.go          # State persistence
├── state/
│   └── last_states.json      # State file (auto-generated)
├── .github/
│   └── workflows/
│       └── rss-monitor.yml   # GitHub Actions workflow
├── feeds.yaml                # Your feed configuration
├── feeds.yaml.example        # Example configuration
├── go.mod
├── go.sum
└── README.md
```

## Troubleshooting

### Notifications Not Received

1. Check that `BARK_DEVICE_KEY` is correctly set in GitHub Secrets
2. Verify your device key by sending a test notification:
   ```bash
   curl "https://api.day.app/YOUR_DEVICE_KEY/Test/Hello"
   ```
3. Check GitHub Actions logs for errors

### State Not Updating

1. Ensure the workflow has write permissions to the repository
2. Check if the branch is protected (may require additional configuration)
3. Review the "Commit state" step in the Actions log

### Feed Not Working

1. Verify the RSS feed URL is accessible
2. Check the feed format (RSS/Atom/JSON Feed)
3. Review the logs for specific error messages

### Workflow Not Running

1. Check that GitHub Actions is enabled in your repository
2. Verify the cron schedule in `.github/workflows/rss-monitor.yml`
3. Note that scheduled workflows may be disabled after 60 days of repository inactivity

## AI Summary Feature

RSS Watcher supports AI-powered summaries using large language models. When configured, push notifications will include AI-generated summaries instead of raw RSS descriptions.

### Environment Variables

Set the following environment variables to enable AI summaries:

```bash
export API_ENDPOINT="https://api.openai.com/v1/chat/completions"
export API_KEY="your-api-key"
export MODEL_NAME="gpt-3.5-turbo"
```

### Supported API Services

- OpenAI API
- Azure OpenAI Service
- Other OpenAI-compatible services (e.g., Ollama, vLLM)

### Features

- ✅ Optional: Not enabled by default, fully backward compatible
- ✅ Auto-fallback: Falls back to original description on API failure
- ✅ Chinese optimized: Optimized for Chinese summaries

For detailed usage, see: [AI Summary Documentation](docs/AI_SUMMARY.md)

## Advanced Usage

### Custom Bark Server

If you're running your own Bark server:

```bash
# Set in GitHub Secrets
BARK_SERVER=https://your-bark-server.com
```

### Adjusting Schedule

Edit `.github/workflows/rss-monitor.yml`:

```yaml
on:
  schedule:
    - cron: '*/15 * * * *'  # Every 15 minutes
    # - cron: '0 * * * *'   # Every hour
    # - cron: '0 */6 * * *' # Every 6 hours
```

### Multiple Device Keys

To send notifications to multiple devices, you can:

1. Use Bark's group feature (same key, different groups)
2. Run multiple instances with different configurations
3. Modify the notifier to support multiple keys

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Security

- Never commit your `BARK_DEVICE_KEY` to the repository
- Use GitHub Secrets for all sensitive data
- Review the state file to ensure no sensitive information is stored

## Credits

- [gofeed](https://github.com/mmcdole/gofeed) - RSS/Atom parser
- [Bark](https://github.com/Finb/Bark) - iOS notification service
- Inspired by various RSS monitoring solutions

## FAQ

**Q: How much does this cost?**  
A: Free! GitHub Actions provides 2,000 minutes/month for free accounts.

**Q: Can I use this for Android?**  
A: You can adapt the notifier to use other services like Telegram, Discord, or Pushover.

**Q: How many feeds can I monitor?**  
A: Limited only by GitHub Actions execution time (typically ~6 hours max). Most users can monitor 50-100 feeds.

**Q: Will this work on a private repository?**  
A: Yes, GitHub Actions works on both public and private repositories.

## Support

If you encounter any issues, please:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review the GitHub Actions logs
3. Open an issue with details about your problem
