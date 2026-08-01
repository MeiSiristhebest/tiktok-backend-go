# Contributing to tiktok-backend-go

We love your input! We want to make contributing to tiktok-backend-go as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Setting Up the Development Environment

**Prerequisites**:
- Go 1.20+
- MySQL 8.0+
- Redis 7+

**Setup steps**:

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/tiktok-backend-go.git`
3. Install dependencies: `go mod tidy`
4. Set up MySQL and Redis (see README.md for port/config defaults)
5. Configure environment variables or update `config/` YAML files
6. Run the service: `go run ./cmd/main.go`

### Branch Naming Convention

- `feature/feature-name`: For new features
- `fix/issue-description`: For bug fixes
- `docs/documentation-update`: For documentation updates
- `refactor/module-name`: For code refactoring
- `perf/description`: For performance-related improvements

### Commit Message Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages:

- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation only changes
- `style:` Changes that do not affect the meaning of the code (formatting, etc)
- `refactor:` A code change that neither fixes a bug nor adds a feature
- `perf:` A code change that improves performance
- `test:` Adding missing tests or correcting existing tests
- `chore:` Changes to the build process or auxiliary tools

### Pull Request Process

1. Create a new branch from `main` following the branch naming convention
2. Make your changes
3. Format the code: `gofmt -w .`
4. Run static analysis: `go vet ./...`
5. Run tests and ensure they pass: `go test ./...`
6. Build the project to ensure it compiles: `go build ./cmd/...`
7. Update documentation if necessary
8. Create a pull request to the `main` branch
9. Wait for review and address any feedback

## Any contributions you make will be under the MIT Software License

In short, when you submit github contributions, you're agreeing to license them under the same terms as the project's license.

## Report bugs using GitHub's issue tracker

We use GitHub issues to track public bugs. Report a bug by opening a new issue; it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Layered Architecture Guidelines

This project follows a layered architecture. When making changes:
- `model/` — data models, do not add business logic here
- `repository/` — data access layer
- `service/` — business logic layer
- `controller/` — HTTP handler layer, keep thin
- `middleware/` — Gin middlewares

## Documentation

- Update README.md with any new features or changes
- Document new API endpoints and request/response schemas

## Questions?

If you have any questions or need help, please open an issue or reach out to the maintainers.

Thank you for contributing!
