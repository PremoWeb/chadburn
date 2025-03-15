# Contributing to Chadburn

Thank you for your interest in contributing to Chadburn! This document provides guidelines and instructions for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.20 or later
- Docker
- Git

### Getting the Code

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/Chadburn.git
   cd Chadburn
   ```

3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/PremoWeb/Chadburn.git
   ```

### Building from Source

```bash
go build -o chadburn .
```

### Running Tests

```bash
go test ./...
```

## Development Workflow

### Creating a Branch

Create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
```

Use a descriptive branch name that reflects the changes you're making.

### Making Changes

1. Make your changes to the codebase
2. Add tests for new functionality
3. Ensure all tests pass
4. Update documentation if necessary

### Commit Messages

Chadburn uses [Conventional Commits](https://www.conventionalcommits.org/) for commit messages. This helps with automatic versioning and changelog generation.

Format:
```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: Code changes that neither fix a bug nor add a feature
- `perf`: Performance improvements
- `test`: Adding or fixing tests
- `chore`: Changes to the build process or auxiliary tools

Example:
```
feat(scheduler): add support for job timeouts

This adds the ability to specify a timeout for jobs, after which they will be terminated.

Closes #123
```

### Submitting a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Go to the [Chadburn repository](https://github.com/PremoWeb/Chadburn) and create a new pull request

3. Fill out the pull request template with details about your changes

4. Wait for a maintainer to review your pull request

## Code Style

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Add comments for exported functions, types, and constants
- Keep functions small and focused
- Use meaningful variable and function names

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass
- Aim for good test coverage

## Documentation

### Code Documentation

- Add comments to exported functions, types, and constants
- Use godoc-compatible comments

### User Documentation

- Update the documentation in the `docs/` directory for user-facing changes
- Keep examples up-to-date
- Use clear, concise language

## Release Process

Chadburn uses semantic versioning and automatic releases based on commit messages:

1. Commits with `fix:` prefix trigger a patch release (1.0.0 -> 1.0.1)
2. Commits with `feat:` prefix trigger a minor release (1.0.0 -> 1.1.0)
3. Commits with `BREAKING CHANGE:` in the footer trigger a major release (1.0.0 -> 2.0.0)

## Getting Help

If you have questions or need help with contributing:

1. Check the [documentation](https://github.com/PremoWeb/Chadburn/tree/main/docs)
2. Open an issue on GitHub
3. Ask in the GitHub Discussions section

## Code of Conduct

Please be respectful and considerate of others when contributing to Chadburn. We strive to maintain a welcoming and inclusive community. 