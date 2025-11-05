# RSS Watcher - Project Summary

## Overview

RSS Watcher is a lightweight, self-hosted RSS monitoring service that runs on GitHub Actions and sends push notifications to iOS devices via Bark.

## Key Features

✅ **Zero-cost operation** - Runs on GitHub Actions free tier  
✅ **iOS notifications** - Push notifications via Bark  
✅ **Multiple feeds** - Monitor unlimited RSS/Atom feeds  
✅ **Smart deduplication** - Track seen items by GUID, link, or title  
✅ **Concurrent processing** - Fast parallel feed fetching  
✅ **Retry logic** - Exponential backoff for failed requests  
✅ **State persistence** - Track feed states in repository  
✅ **Aggregated notifications** - Combine multiple updates into one notification  
✅ **Easy deployment** - Fork, configure, done!  

## Architecture

### Components

1. **Fetcher** (`internal/fetcher/`)
   - HTTP client with timeouts and retry logic
   - Exponential backoff for failures
   - Custom User-Agent for RSS feeds

2. **Parser** (`internal/parser/`)
   - Based on gofeed library
   - Supports RSS, Atom, and JSON Feed
   - Extracts title, description, link, published date

3. **Deduper** (`internal/deduper/`)
   - Tracks last seen item per feed
   - Configurable deduplication key (guid/link/title)
   - Efficiently identifies new items

4. **Notifier** (`internal/notifier/`)
   - Bark API integration
   - Individual or aggregated notifications
   - URL encoding for special characters
   - Group support for organization

5. **State Manager** (`internal/state/`)
   - JSON-based state storage
   - Atomic writes (write to .tmp, then rename)
   - Thread-safe operations with mutex

6. **Config Manager** (`internal/config/`)
   - YAML configuration parsing
   - Feed-specific settings
   - Validation

### Data Flow

```
GitHub Actions (cron) 
    ↓
Checkout Repository
    ↓
Build Go Binary
    ↓
Load Config (feeds.yaml)
    ↓
Load State (last_states.json)
    ↓
Fetch Feeds (concurrent, max 8)
    ↓
Parse RSS/Atom
    ↓
Deduplicate Items
    ↓
Send Bark Notifications
    ↓
Update State
    ↓
Commit State to Repository
```

## File Structure

```
.
├── cmd/rsswatcher/           # Main application
│   └── main.go              # Entry point with CLI flags
├── internal/                # Internal packages
│   ├── config/              # Configuration handling
│   ├── deduper/             # Deduplication logic
│   ├── fetcher/             # HTTP fetching
│   ├── notifier/            # Bark notifications
│   ├── parser/              # RSS parsing
│   └── state/               # State persistence
├── .github/workflows/       # GitHub Actions
│   ├── rss-monitor.yml     # Main scheduled workflow
│   └── test.yml            # CI testing workflow
├── state/                  # State storage
│   └── last_states.json    # Feed states (auto-managed)
├── feeds.yaml              # User configuration
├── feeds.yaml.example      # Example configuration
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile               # Build automation
├── test-local.sh          # Local testing script
├── README.md              # Main documentation
├── QUICKSTART.md          # Quick start guide
├── DEPLOYMENT.md          # Deployment guide
├── CONTRIBUTING.md        # Contribution guidelines
└── LICENSE                # MIT License
```

## Configuration

### feeds.yaml

```yaml
feeds:
  - id: unique-feed-id          # Required: Unique identifier
    name: Display Name           # Required: Shown in notifications
    url: https://example.com/rss # Required: Feed URL
    notify: true                 # Optional: Enable notifications
    dedupe_key: guid             # Optional: guid|link|title
    aggregate: false             # Optional: Aggregate notifications
    aggregate_window_minutes: 30 # Optional: Aggregation window
```

### Environment Variables (GitHub Secrets)

- `BARK_DEVICE_KEY` (Required) - Your Bark device key
- `BARK_SERVER` (Optional) - Custom Bark server URL

## Development

### Prerequisites

- Go 1.22+
- Git

### Commands

```bash
# Build
make build

# Test
make test

# Run locally
./rsswatcher --config feeds.yaml --state state/last_states.json

# Full validation
make all

# Clean
make clean
```

### Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -v -race -coverprofile=coverage.txt ./...

# Test locally with Bark
./test-local.sh YOUR_DEVICE_KEY
```

## Deployment

1. Fork repository
2. Get Bark device key from iOS app
3. Add `BARK_DEVICE_KEY` to GitHub Secrets
4. Edit `feeds.yaml` with your RSS feeds
5. Commit and push
6. Workflow runs automatically every 30 minutes

See [DEPLOYMENT.md](../DEPLOYMENT.md) for detailed instructions.

## Technical Decisions

### Why Go?
- Single static binary
- Excellent concurrency
- Fast compilation
- No runtime dependencies
- Great HTTP/parsing libraries

### Why GitHub Actions?
- Free tier sufficient for most users
- Built-in secrets management
- Easy to audit and modify
- No server maintenance
- Git-based state storage

### Why Bark?
- Simple iOS notifications
- No account required (for default server)
- Self-hosting option available
- Free and open source
- REST API (no SDK needed)

### Why Repository State Storage?
- Simple and auditable
- No external dependencies
- Free (included with repository)
- Version controlled
- Works with GitHub Actions

## Performance

### Typical Metrics
- Feeds: Up to 100 feeds comfortably
- Execution time: 30-120 seconds (depends on feed count)
- Concurrency: Max 8 concurrent fetches
- Timeout: 15 seconds per feed, 30 seconds overall
- Retry: Up to 2 retries with exponential backoff

### Resource Usage
- Binary size: ~11 MB
- Memory: <50 MB during execution
- GitHub Actions minutes: ~2 minutes per run
- Monthly cost: $0 (free tier)

## Monitoring

### Check Status
- GitHub Actions tab shows all runs
- View logs for each execution
- State file shows last processed items

### Error Handling
- Failed fetches are logged but don't stop other feeds
- Bark failures are logged
- State updates are atomic (all-or-nothing)

## Future Enhancements

Potential features for community contributions:

- [ ] Support for more notification services (Telegram, Discord, etc.)
- [ ] Web UI for configuration
- [ ] Email notifications
- [ ] RSS feed validation and health checks
- [ ] Custom notification templates
- [ ] Filter rules (regex, keywords)
- [ ] Feed categories and grouping
- [ ] Statistics and analytics
- [ ] Multiple Bark devices
- [ ] Webhook support
- [ ] Docker support for self-hosting

## Security

- Secrets stored in GitHub Secrets (encrypted)
- State file contains no sensitive data
- HTTPS for all external requests
- No external database or storage
- Minimal attack surface

## License

MIT License - See [LICENSE](../LICENSE)

## Support

- Documentation: [README.md](../README.md)
- Quick Start: [QUICKSTART.md](../QUICKSTART.md)
- Deployment: [DEPLOYMENT.md](../DEPLOYMENT.md)
- Issues: GitHub Issues
- Contributions: [CONTRIBUTING.md](../CONTRIBUTING.md)

## Credits

- **gofeed**: RSS/Atom parsing
- **Bark**: iOS notifications
- **GitHub Actions**: CI/CD platform
- **Go**: Programming language

---

Built with ❤️ for the RSS community
