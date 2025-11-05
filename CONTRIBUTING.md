# Contributing to RSS Watcher

Thank you for your interest in contributing to RSS Watcher! This document provides guidelines and instructions for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/rsswatcher.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `go test ./...`
6. Build the project: `go build -o rsswatcher ./cmd/rsswatcher`
7. Commit your changes: `git commit -am 'Add some feature'`
8. Push to the branch: `git push origin feature/your-feature-name`
9. Create a Pull Request

## Development Guidelines

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Keep functions small and focused
- Write clear, descriptive variable and function names

### Testing

- Write tests for new functionality
- Ensure all tests pass before submitting a PR
- Aim for good test coverage
- Use table-driven tests where appropriate

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in the imperative mood (e.g., "Add", "Fix", "Update")
- Keep the first line under 72 characters
- Add a blank line followed by a detailed description if needed

Example:
```
Add support for JSON Feed format

- Implement JSON Feed parser
- Add tests for JSON Feed parsing
- Update documentation with JSON Feed examples
```

## Project Structure

```
.
├── cmd/rsswatcher/      # Main application
├── internal/            # Internal packages
│   ├── config/         # Configuration handling
│   ├── deduper/        # Deduplication logic
│   ├── fetcher/        # HTTP fetching with retry
│   ├── notifier/       # Notification services
│   ├── parser/         # Feed parsing
│   └── state/          # State persistence
└── .github/workflows/  # GitHub Actions workflows
```

## Adding New Features

### Adding a New Notification Service

1. Create a new file in `internal/notifier/`
2. Implement the notification interface
3. Update the main application to support the new notifier
4. Add tests for the new notifier
5. Update the README with usage instructions

### Adding New Feed Formats

1. Update the parser in `internal/parser/`
2. Add tests with sample feed data
3. Update documentation

## Reporting Issues

When reporting issues, please include:

- A clear description of the problem
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Your environment (OS, Go version, etc.)
- Relevant logs or error messages

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update tests to cover your changes
3. Ensure all tests pass
4. Update documentation as needed
5. The PR will be merged once reviewed and approved

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment for all contributors

## Questions?

If you have questions, feel free to:
- Open an issue for discussion
- Check existing issues and PRs for similar topics

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
