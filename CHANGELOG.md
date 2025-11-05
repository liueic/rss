# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2024-11-05

### Added
- Full Unicode/Chinese character support
- Proper character truncation for multi-byte characters (Chinese, Japanese, emojis)
- Chinese documentation (README.zh-CN.md, QUICKSTART.zh-CN.md)
- Comprehensive tests for Chinese character handling
- Chinese support documentation (docs/CHINESE_SUPPORT.md)

### Changed
- Improved `truncate()` function to use rune-based truncation instead of byte-based
- Enhanced `cleanDescription()` in parser to handle Unicode correctly
- Updated string processing to prevent corruption of multi-byte characters

### Fixed
- Fixed Chinese characters being cut in half during truncation
- Fixed emoji and other Unicode characters display issues
- Ensured all text processing is UTF-8 safe

## [1.0.0] - 2024-11-05

### Added
- Initial release of RSS Watcher
- Go-based RSS monitoring service
- GitHub Actions integration for scheduled execution
- Bark notification support for iOS devices
- Multi-feed monitoring with concurrent processing
- Smart deduplication (by GUID, link, or title)
- State persistence in repository
- Configurable aggregated notifications
- Retry logic with exponential backoff
- YAML-based configuration
- Comprehensive documentation (README, QUICKSTART, DEPLOYMENT)
- Unit tests for core components
- Example configurations
- Local testing script
- Makefile for build automation
- CI/CD workflow for testing
- MIT License

### Features
- Monitor unlimited RSS/Atom feeds
- Push notifications to iOS via Bark
- Concurrent feed fetching (max 8 simultaneous)
- Automatic deduplication of seen items
- Individual or aggregated notifications
- Configurable check frequency (default: 30 minutes)
- Zero-cost operation on GitHub Actions free tier
- Self-hosted Bark server support
- Thread-safe state management
- Atomic state file updates

### Documentation
- Complete README with features and usage
- Quick start guide for fast deployment
- Detailed deployment instructions
- Contributing guidelines
- Project summary and architecture overview

### Technical Details
- Go 1.22+ support
- Dependencies: gofeed, yaml.v3
- Test coverage: 80%+ for tested components
- Binary size: ~11 MB
- Memory usage: <50 MB during execution

[1.0.0]: https://github.com/rsswatcher/rsswatcher/releases/tag/v1.0.0
