# Contributing to Wafer

Thank you for your interest in contributing to wafer! This document provides guidelines and information for contributors, following the structured approach used by the mochilang/mochi project.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Set up the development environment**
4. **Make your changes**
5. **Test your changes**
6. **Submit a pull request**

## 🛠️ Development Setup

### Prerequisites

- Go 1.23 or later
- Make
- Docker (optional)
- Git

### Environment Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wafer.git
cd wafer

# Set up development environment
make dev-setup

# Verify setup
make dev-check

# Run tests to ensure everything works
make test
```

## 📋 Contribution Guidelines

### Code Style

- **Go formatting**: Use `gofmt` (run `make fmt`)
- **Linting**: Pass `golangci-lint` checks (run `make lint`)
- **Comments**: Document public functions and complex logic
- **Error handling**: Always handle errors appropriately
- **Testing**: Write tests for new functionality

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(chunker): add semantic chunking support
fix(embedder): handle timeout errors gracefully
docs(readme): update installation instructions
test(golden): add regression tests for JSONL output
```

### Pull Request Process

1. **Create a feature branch** from `main`
2. **Make your changes** with appropriate tests
3. **Update documentation** if needed
4. **Run the full test suite**
5. **Submit a pull request** with a clear description

#### PR Requirements

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] Benchmarks don't regress significantly
- [ ] Golden tests pass (if applicable)

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test types
make test-unit           # Unit tests
make test-integration    # Integration tests
make test-golden         # Golden file tests

# Run benchmarks
make bench
```

### Writing Tests

#### Unit Tests
- Test individual components in isolation
- Use table-driven tests for multiple scenarios
- Mock external dependencies

#### Integration Tests
- Test complete workflows
- Use temporary directories for file operations
- Clean up resources in test teardown

#### Golden File Tests
- Add test input files to `tests/golden/`
- Run `make update-golden` to generate expected output
- Commit both input and golden files

#### Benchmarks
- Add benchmark functions to `tests/benchmark_test.go`
- Use `testing.B` for performance measurements
- Include memory allocation measurements

### Test Coverage

Maintain high test coverage:
- **Target**: >90% overall coverage
- **Critical paths**: 100% coverage
- **New features**: Must include tests

## 📊 Performance

### Benchmark Guidelines

- **Baseline**: Set performance baselines with `make bench-baseline`
- **Regression**: Avoid performance regressions >10%
- **Profiling**: Use `make bench-profile` for optimization
- **Documentation**: Update BENCHMARK.md for significant changes

### Performance Considerations

- **Memory efficiency**: Minimize allocations
- **CPU usage**: Optimize hot paths
- **I/O patterns**: Use efficient file operations
- **Concurrency**: Consider parallel processing opportunities

## 📝 Documentation

### Documentation Requirements

- **Code comments**: Document public APIs
- **README updates**: Update for new features
- **Guides**: Add usage examples
- **Changelog**: Update CHANGELOG.md

### Documentation Structure

```
docs/
├── guides/          # User guides
├── features/        # Feature documentation
├── tutorials/       # Step-by-step tutorials
├── examples/        # Code examples
└── api/            # API reference
```

## 🐛 Bug Reports

### Before Reporting

1. **Search existing issues** for duplicates
2. **Test with latest version**
3. **Reproduce with minimal example**
4. **Check documentation** for expected behavior

### Bug Report Template

```markdown
**Description**
Brief description of the bug

**Steps to Reproduce**
1. Step one
2. Step two
3. Step three

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.23]
- Wafer version: [e.g., v0.1.0]
- Ollama version: [e.g., 0.1.17]

**Additional Context**
Any other relevant information
```

## 💡 Feature Requests

### Before Requesting

1. **Check the roadmap** (ROADMAP.md)
2. **Search existing issues** for similar requests
3. **Consider the scope** and impact
4. **Think about implementation** approach

### Feature Request Template

```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this feature work?

**Alternatives Considered**
Other approaches you've considered

**Additional Context**
Any other relevant information
```

## 🏗️ Architecture

### Project Structure

```
wafer/
├── cmd/wafer/          # CLI entrypoint
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   └── ingest/         # Core processing logic
├── golden/             # Golden file testing framework
├── tests/              # Test files
├── docs/               # Documentation
└── bench/              # Benchmark results
```

### Design Principles

- **Modularity**: Keep components loosely coupled
- **Testability**: Design for easy testing
- **Performance**: Optimize for speed and memory
- **Reliability**: Handle errors gracefully
- **Usability**: Provide clear interfaces and messages

## 🔄 Release Process

### Version Management

- **Semantic versioning**: MAJOR.MINOR.PATCH
- **Version file**: Update VERSION file
- **Changelog**: Update CHANGELOG.md
- **Tags**: Create git tags for releases

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Version bumped
- [ ] Performance benchmarks reviewed
- [ ] Security review completed

## 🤝 Community

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions
- **Documentation**: Improvements and clarifications

### Code of Conduct

We follow the [Contributor Covenant](https://www.contributor-covenant.org/) code of conduct. Please be respectful and inclusive in all interactions.

### Recognition

Contributors are recognized in:
- **CHANGELOG.md**: Major contributions
- **README.md**: Acknowledgments section
- **GitHub**: Contributor graphs and statistics

## 📚 Resources

### Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Testing in Go](https://golang.org/doc/tutorial/add-a-test)

### Project Resources

- [Technical Specification](SPEC.md)
- [Roadmap](ROADMAP.md)
- [Benchmark Results](BENCHMARK.md)
- [Usage Guide](docs/guides/usage.md)

## ❓ Getting Help

### Before Asking

1. **Read the documentation**
2. **Search existing issues**
3. **Check the FAQ** (if available)
4. **Try the latest version**

### How to Ask

- **Be specific** about the problem
- **Provide context** and environment details
- **Include relevant code** or configuration
- **Show what you've tried**

### Response Time

- **Bug reports**: 1-3 business days
- **Feature requests**: 1-2 weeks
- **Pull requests**: 1-5 business days
- **Questions**: 1-3 business days

---

Thank you for contributing to wafer! Your contributions help make this project better for everyone. 🙏
