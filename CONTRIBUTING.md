# Contributing to Common Library

## Development Process

### Prerequisites
- Go 1.24 or later
- Git
- Make (optional, for using Makefile commands)

### Getting Started
1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/common.git`
3. Add upstream: `git remote add upstream https://github.com/kopexa-grc/common.git`

### Development Workflow
1. Create a new branch from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following our coding standards
3. Write or update tests as needed
4. Run tests locally: `go test ./...`
5. Commit your changes using conventional commits

### Handling Dependabot Pull Requests

We use Dependabot to keep our dependencies up to date. When reviewing Dependabot PRs:

1. Check the changelog/release notes of the updated dependency
2. Verify that the update is compatible with our codebase
3. Run the test suite to ensure nothing breaks
4. Pay special attention to security-related updates
5. If the update is a major version, coordinate with the team

### Commit Guidelines
We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or modifying tests
- `chore:` - Maintenance tasks
- `security:` - Security-related changes

Example:
```
feat: add new validation function for email addresses
```

### Pull Request Process
1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass
4. Update the CHANGELOG.md (this will be done automatically by release-please)
5. Request review from at least one maintainer

### Code Review Guidelines
- All PRs require at least one approval
- Address all review comments
- Keep PRs focused and manageable in size
- Use meaningful commit messages
- Follow Go best practices and idioms

## Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Follow our linting rules
- Write clear and concise comments
- Use meaningful variable and function names

## Testing
- Write unit tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies

## Documentation
- Update README.md for significant changes
- Document all exported functions and types
- Include examples where helpful
- Keep documentation up to date with code changes 