# Contributing to Prompt Optimizer

First off, thank you for considering contributing to Prompt Optimizer! It's people like you that make this tool better for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Describe the behavior you observed and what you expected**
- **Include your environment details** (OS, Claude Code version, etc.)

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- **A clear and descriptive title**
- **A detailed description of the proposed enhancement**
- **Explain why this enhancement would be useful**

### Pull Requests

1. **Fork the repo** and create your branch from `main`
2. **Follow the existing code style** - check existing files for patterns
3. **Test your changes** - run `go test ./...` for backend, `npm test` for frontend
4. **Update documentation** if needed
5. **Write a clear commit message**

## Development Setup

### Prerequisites

- Go 1.24+
- Node.js 18+
- Wails CLI v2

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/prompt-optimizer-skill.git
cd prompt-optimizer-skill

# Backend tests
cd wails-app
go test -v ./...

# Frontend development
cd frontend
npm install
npm run dev
```

### Project Structure

```
prompt-optimizer-skill/
├── skills/           # Skill definition files (user-facing)
├── wails-app/        # Desktop application source
│   ├── main.go       # CLI + Wails initialization
│   ├── app.go        # Go bindings for frontend
│   └── frontend/     # Pure HTML/CSS/JS frontend
└── .github/          # CI/CD workflows
```

## Code Review Process

1. A maintainer will review your PR
2. Changes may be requested - this is normal!
3. Once approved, your PR will be merged

## Questions?

Feel free to open an issue with the "question" label.

---

Thank you for contributing!
