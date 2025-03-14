# Contributing to Chadburn

Thank you for your interest in contributing to Chadburn! This document provides guidelines for contributing to the project.

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) for our commit messages. This enables automatic versioning and changelog generation.

### Commit Message Format

Each commit message consists of a **header**, a **body**, and a **footer**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory, while the **body** and **footer** are optional.

#### Type

The type must be one of the following:

- **feat**: A new feature (triggers a minor version bump)
- **fix**: A bug fix (triggers a patch version bump)
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries

#### Breaking Changes

Commits with breaking changes should include `BREAKING CHANGE:` in the footer or append a `!` after the type/scope. This will trigger a major version bump.

Example:
```
feat!: remove support for older Docker versions

BREAKING CHANGE: Docker versions below 19.03 are no longer supported.
```

#### Examples

```
feat(scheduler): add support for timezone specification
```

```
fix(logging): correct timestamp format in log output
```

```
docs: update README with new configuration options
```

```
chore(deps): update dependencies
```

### Git Hooks

This repository uses Git hooks to enforce commit message format and code quality. The hooks are managed by Husky and commitlint.

#### Setup

After cloning the repository, run:

```bash
npm install
```

This will set up the Git hooks automatically.

#### Commit Message Validation

When you commit changes, the commit-msg hook will validate your commit message against the Conventional Commits format. If the message doesn't follow the format, the commit will be rejected.

#### Pre-commit Checks

The pre-commit hook runs various checks before allowing a commit, such as:
- Code formatting
- Linting
- Tests

If any of these checks fail, the commit will be rejected.

## Pull Request Process

1. Ensure your code adheres to the project's coding standards.
2. Update the README.md with details of changes to the interface, if applicable.
3. The versioning scheme we use is [SemVer](http://semver.org/).
4. Your pull request may be merged once it passes all tests and has been reviewed by maintainers.

## Code of Conduct

Please be respectful and inclusive in your interactions with the project community. We aim to foster an open and welcoming environment. 