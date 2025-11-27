# H.A.T. Stack Bootstrap

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![HTMX](https://img.shields.io/badge/HTMX-2.0-orange)](https://htmx.org)
[![Templ](https://img.shields.io/badge/Templ-Latest-green)](https://templ.guide)

A production-ready Go web application template built with the **H.A.T. Stack**: **HTMX**, **Alpine.js**, and **Templ**.

This bootstrap provides everything you need to start building modern web applications with Go, featuring a clean architecture, type-safe templates, and live-reload development.

## 🚀 Quick Start

**Want to get started immediately?** Check out the [Quick Start Guide](Docs/QuickStart.md) for a 5-minute setup!

## Table of Contents

- [What is the H.A.T. Stack?](#what-is-the-hat-stack)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Option 1: Download Latest Release](#option-1-download-latest-release-recommended-for-new-projects)
  - [Option 2: Clone the Repository](#option-2-clone-the-repository)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Development](#development)
- [Customizing for Your Project](#customizing-for-your-project)
- [Contributing](#contributing)
- [Why H.A.T. Stack?](#why-hat-stack)

## What is the H.A.T. Stack?

- **HTMX** - Server interactions without writing JavaScript
- **Alpine.js** - Lightweight client-side reactivity (when you need it)
- **Templ** - Type-safe Go templates with compile-time checking

Combined with Go's performance and simplicity, the H.A.T. Stack enables rapid development of modern web applications with minimal frontend complexity.

## Features

- ✅ **Dual Architecture**: Support for both JSON API endpoints (`/api/*`) and server-rendered HTML pages
- ✅ **Repository Pattern**: Clean data access layer with Google Cloud Datastore
- ✅ **Service Layer**: Business logic separation with proper dependency injection
- ✅ **Type-Safe Templates**: Templ provides compile-time type safety for HTML templates
- ✅ **Live Reload**: Air for automatic rebuilding during development
- ✅ **Modern Frontend**: HTMX for dynamic interactions, Alpine.js for reactivity, TailwindCSS for styling
- ✅ **Cross-Platform**: Works on Windows, Linux, and macOS

## Tech Stack

### Backend
- **Go** - Primary language
- **Gin** - Web framework
- **Templ** - Type-safe Go templates
- **Zerolog** - Structured logging
- **Google Cloud Datastore** - NoSQL database

### Frontend
- **HTMX** - Server interactions without JavaScript
- **Alpine.js** - Lightweight client-side reactivity
- **TailwindCSS** - Utility-first CSS framework

## Getting Started

### Prerequisites

- **Go 1.23 or higher** - [Download Go](https://golang.org/dl/)
- **Google Cloud account** - For Datastore (optional, can be replaced with another database)

### Option 1: Download Latest Release (Recommended for New Projects)

Perfect if you want to start your own repository from scratch:

1. **Download the latest release:**
   ```bash
   # Download and extract
   curl -L https://github.com/runtime-dynamics/hatstack/archive/refs/heads/main.zip -o hatstack.zip
   unzip hatstack.zip
   cd hatstack-main
   
   # Or use wget
   wget https://github.com/runtime-dynamics/hatstack/archive/refs/heads/main.zip
   unzip main.zip
   cd hatstack-main
   ```

2. **Initialize your own Git repository:**
   ```bash
   rm -rf .git  # Remove the bootstrap's git history
   git init
   git add .
   git commit -m "Initial commit from H.A.T. Stack bootstrap"
   ```

3. **Update the module name** in `go.mod`:
   ```go
   module github.com/yourusername/yourproject
   ```

4. **Run setup:**
   
   **Windows:**
   ```cmd
   setup.bat
   ```
   
   **Linux/Mac:**
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```

### Option 2: Clone the Repository

Best if you want to contribute back or stay updated with bootstrap improvements:

```bash
# Clone the repository
git clone https://github.com/runtime-dynamics/hatstack.git
cd hatstack

# Run setup
# Windows:
setup.bat

# Linux/Mac:
chmod +x setup.sh
./setup.sh
```

### What Setup Does

The setup script will:
1. ✅ Install **Air** (live-reload tool) if not present
2. ✅ Install **Templ** (template generator) if not present
3. ✅ Create `.env` file with template (you'll need to configure it)
4. ✅ Set up platform-specific Air configuration
5. ✅ Make scripts executable (Linux/Mac)

### Configuration

Edit the `.env` file with your settings:

```env
# Google Cloud Configuration
DATASTORE_NAME=your-datastore-name
PUBSUB_TOPIC=your-pubsub-topic
PUBSUB_SUBSCRIPTION=your-pubsub-subscription

# Frontend Configuration
FRONTEND_ENDPOINT=http://local.nitecon.net:8080

# Google Cloud Authentication
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
# OR
GOOGLE_PROJECT_ID=your-project-id

# Development Settings
DEBUG=true
PORT=8080
```

### Start Development Server

```bash
air
```

**Or use the quick launcher (Linux/Mac):**
```bash
./dev.sh
```

The server will:
- ✅ Generate Templ templates automatically
- ✅ Load environment variables from `.env`
- ✅ Build and start the application
- ✅ Rebuild on file changes (`.go`, `.templ`, `.html`)
- ✅ Hot-reload in your browser

**Access the application:** `http://localhost:8080`

You should see the H.A.T. Stack welcome page! 🎉

### Next Steps

1. **Customize the homepage** - Edit `views/pages/home.templ`
2. **Add your first API endpoint** - See `web/api/routes.go`
3. **Create a new page** - Add a handler in `web/app/` and template in `views/pages/`
4. **Read the coding guidelines** - Check `Docs/CodingGuidelines.md` for best practices

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── config/                # Configuration management
├── data/                  # Data layer (repositories)
├── services/              # Business logic layer
├── web/                   # Web layer
│   ├── api/              # JSON API handlers (/api/*)
│   └── app/              # HTML page handlers
├── views/                 # Templ templates
│   ├── components/       # Reusable UI components
│   ├── layouts/          # Page layouts
│   └── pages/            # Page templates
├── static/               # Static assets (CSS, JS, images)
├── scripts/              # Development scripts
└── Docs/                 # Documentation
    └── CodingGuidelines.md
```

## Architecture

### Dual Handler Pattern

**API Handlers** (`web/api/`):
- Handle `/api/*` routes
- Return JSON responses
- Used for AJAX/fetch requests

**App Handlers** (`web/app/`):
- Handle all other routes
- Return HTML via Templ templates
- Server-side rendering

### Layer Separation

```
Handler → Service → Repository → Datastore
```

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and orchestration
- **Repositories**: Data access and persistence

## Development

### Generate Templ Templates

```bash
templ generate
```

Air runs this automatically before each build.

### Build

```bash
go build -o ./tmp/main ./cmd
```

### Run Tests

```bash
go test ./...
```

## Customizing for Your Project

After downloading/cloning the bootstrap, you'll want to customize it:

### 1. Update Module Name

Edit `go.mod` and change the module path:
```go
module github.com/yourusername/yourproject
```

Then update all imports in your Go files to match.

### 2. Update Branding

- Edit `views/pages/home.templ` - Update the homepage content
- Replace `static/images/logo-square_128.png` with your logo
- Update `static/favicon.ico` with your favicon

### 3. Configure Your Database

The bootstrap uses Google Cloud Datastore by default. To use a different database:

1. Update the repository implementations in `data/`
2. Modify `data/data.go` to initialize your database client
3. Update the `.env` configuration

### 4. Remove Unused Features

The bootstrap includes examples for common patterns. Remove what you don't need:
- Example repositories in `data/`
- Example services in `services/`
- Example handlers in `web/api/` and `web/app/`

## Deployment

This project includes a production-ready `Dockerfile` that can be deployed directly to **Google Cloud Run** or any container platform.

### Docker Build

```bash
# Build the Docker image
docker build -t hatstack-app .

# Run locally with Docker
docker run -p 8080:8080 --env-file .env hatstack-app

# Test the container
curl http://localhost:8080
```

### Deploy to Google Cloud Run

```bash
# Deploy directly from source
gcloud run deploy hatstack-app --source .

# Or build and deploy
docker build -t gcr.io/YOUR-PROJECT-ID/hatstack-app .
docker push gcr.io/YOUR-PROJECT-ID/hatstack-app
gcloud run deploy hatstack-app --image gcr.io/YOUR-PROJECT-ID/hatstack-app
```

### Environment Variables

Set these environment variables in your deployment platform:

- `DATASTORE_NAME` - Your Datastore database name
- `GOOGLE_PROJECT_ID` - Your Google Cloud project ID
- `FRONTEND_ENDPOINT` - Your application's public URL
- `PORT` - Port to listen on (default: 8080)

### Container Features

The Dockerfile is optimized for production:
- ✅ Multi-stage build for minimal image size
- ✅ Templ templates pre-generated
- ✅ Static assets included
- ✅ Non-root user for security
- ✅ Health check endpoint ready
- ✅ Cloud Run compatible

## Testing

The bootstrap includes comprehensive tests for all components.

**Run all tests:**
```bash
# Linux/Mac
make test

# Windows
test.bat test
```

**Run with coverage:**
```bash
# Linux/Mac
make test-coverage

# Windows
test.bat coverage
```

**Run benchmarks:**
```bash
# Linux/Mac
make bench

# Windows
test.bat bench
```

See **[Testing Guide](Docs/Testing.md)** for detailed testing documentation.

## Documentation

- **[Quick Start Guide](Docs/QuickStart.md)** - Get started in 5 minutes
- **[Coding Guidelines](Docs/CodingGuidelines.md)** - Development standards and patterns
- **[Testing Guide](Docs/Testing.md)** - Testing strategy and best practices
- **[Scripts Documentation](Docs/Scripts.md)** - Development scripts reference

## Contributing

We welcome contributions to improve the H.A.T. Stack bootstrap!

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines

- Follow the existing code style and patterns
- Update documentation for any new features
- Keep the bootstrap generic and reusable
- Test on both Windows and Linux/Mac if possible

## Community & Support

- **GitHub Issues** - [Report bugs or request features](https://github.com/runtime-dynamics/hatstack/issues)
- **Discussions** - [Ask questions and share ideas](https://github.com/runtime-dynamics/hatstack/discussions)

## Why H.A.T. Stack?

The H.A.T. Stack offers a refreshing alternative to heavy JavaScript frameworks:

- ✅ **Less JavaScript** - HTMX handles most interactions server-side
- ✅ **Type Safety** - Templ provides compile-time template checking
- ✅ **Fast Development** - Live reload and minimal build steps
- ✅ **Better Performance** - Server-side rendering with minimal client-side JS
- ✅ **Easier Debugging** - Server-side logic is easier to trace and test
- ✅ **SEO Friendly** - Full HTML rendering on the server

Perfect for:
- Internal tools and dashboards
- Content-heavy websites
- CRUD applications
- MVPs and prototypes
- Teams that prefer backend development

## License

MIT License - feel free to use this bootstrap for any project!

Copyright © 2025 Runtime Dynamics LLC. All rights reserved.
