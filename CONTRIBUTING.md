# Contributing to H.A.T. Stack Bootstrap

Thank you for your interest in contributing to the H.A.T. Stack Bootstrap! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Keep discussions professional

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/runtime-dynamics/hatstack/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Your environment (OS, Go version, etc.)
   - Screenshots if applicable

### Suggesting Features

1. Check [Discussions](https://github.com/runtime-dynamics/hatstack/discussions) for similar ideas
2. Create a new discussion or issue explaining:
   - The problem you're trying to solve
   - Your proposed solution
   - Why it would benefit the bootstrap
   - Any alternatives you've considered

### Submitting Changes

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/hatstack.git
   cd hatstack
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Follow the existing code style
   - Keep changes focused and atomic
   - Write clear commit messages

4. **Test your changes**
   ```bash
   # Run tests
   go test ./...
   
   # Test on your platform
   ./setup.sh  # or setup.bat
   air
   ```

5. **Update documentation**
   - Update README.md if needed
   - Update documentation in Docs/ for architectural changes
   - Add comments for complex code

6. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add amazing feature"
   ```
   
   Use conventional commit messages:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes (formatting, etc.)
   - `refactor:` - Code refactoring
   - `test:` - Adding or updating tests
   - `chore:` - Maintenance tasks

7. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create a Pull Request**
   - Go to the original repository
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template
   - Link any related issues

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions small and focused
- Write descriptive variable names
- Add comments for complex logic

### Architecture Principles

- **Keep it generic** - The bootstrap should work for any project
- **Maintain separation of concerns** - Handler → Service → Repository
- **Type safety** - Use Templ for templates, avoid `interface{}`
- **Cross-platform** - Test on Windows and Linux/Mac when possible

### What to Contribute

**Good contributions:**
- Bug fixes
- Performance improvements
- Better error handling
- Cross-platform compatibility fixes
- Documentation improvements
- Generic utility functions
- Better examples in CodingGuidelines.md

**Avoid:**
- Project-specific features
- Breaking changes without discussion
- Large refactors without prior agreement
- Dependencies that aren't widely used

### Testing

- Test your changes locally
- Ensure the setup scripts work
- Verify Air live-reload works
- Test on both Windows and Unix if possible
- Check that Templ generation works

### Documentation

When adding features:
- Update README.md if it affects getting started
- Update Docs/CodingGuidelines.md for architectural patterns
- Add inline comments for complex code
- Update Docs/Scripts.md for script changes
- Update Docs/Testing.md for new testing patterns

## Questions?

- Open a [Discussion](https://github.com/runtime-dynamics/hatstack/discussions) for questions
- Check existing issues and discussions first
- Be patient - maintainers are volunteers

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions

Thank you for helping make the H.A.T. Stack Bootstrap better! 🎉
