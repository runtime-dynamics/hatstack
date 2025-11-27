# Quick Start Guide

Get up and running with the H.A.T. Stack Bootstrap in 5 minutes!

## 🚀 Fast Track (Copy & Paste)

### For Linux/Mac:

```bash
# Download and extract
curl -L https://github.com/runtime-dynamics/hatstack/archive/refs/heads/main.zip -o hatstack.zip
unzip hatstack.zip && cd hatstack-main

# Setup (installs Air, Templ, creates .env)
chmod +x setup.sh && ./setup.sh

# Edit your .env file with your settings
nano .env

# Start development server
air
```

### For Windows:

```cmd
REM Download from: https://github.com/runtime-dynamics/hatstack/archive/refs/heads/main.zip
REM Extract the zip file
cd hatstack-main

REM Run setup
setup.bat

REM Edit .env file with your settings
notepad .env

REM Start development server
air
```

## 📋 What You Need

- **Go 1.24+** installed ([Download](https://golang.org/dl/))
- **Google Cloud account** (for Datastore) or plan to use another database

## 🎯 What Happens During Setup

1. ✅ Checks if Go is installed
2. ✅ Installs Air (live-reload tool)
3. ✅ Installs Templ (template generator)
4. ✅ Creates `.env` file with template
5. ✅ Configures Air for your platform

## ⚙️ Configure Your Environment

Edit the `.env` file created by setup:

```env
# Required: Your Google Cloud settings
DATASTORE_NAME=your-datastore-name
PUBSUB_TOPIC=your-pubsub-topic
PUBSUB_SUBSCRIPTION=your-pubsub-subscription

# Required: Frontend URL
FRONTEND_ENDPOINT=http://localhost:8080

# Required: Authentication (choose one)
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
# OR
GOOGLE_PROJECT_ID=your-project-id

# Optional: Development settings
DEBUG=true
PORT=8080
```

## 🎨 Customize for Your Project

### 1. Update Module Name

Edit `go.mod`:
```go
module github.com/yourusername/yourproject
```

### 2. Update Imports

Find and replace `runtime-dynamics` with your module name in all `.go` files.

### 3. Start Building!

- **Homepage**: Edit `views/pages/home.templ`
- **API Endpoints**: Add handlers in `web/api/`
- **Web Pages**: Add handlers in `web/app/`
- **Business Logic**: Add services in `services/`
- **Data Access**: Add repositories in `data/`

## 🔥 Development Workflow

```bash
# Start dev server (auto-reloads on changes)
air

# In another terminal, make changes to:
# - views/pages/home.templ
# - web/app/home.go
# - Any .go file

# Air will automatically:
# 1. Run templ generate
# 2. Rebuild the app
# 3. Restart the server
# 4. Refresh your browser
```

## 📚 Next Steps

1. Read the [Coding Guidelines](CodingGuidelines.md)
2. Explore the [project structure](../README.md#project-structure)
3. Check out example handlers in `web/api/` and `web/app/`
4. Learn about the [architecture](../README.md#architecture)
5. Review [deployment options](../README.md#deployment) including Docker and Cloud Run

## 🐳 Docker & Cloud Run

This project includes a `Dockerfile` that can be deployed directly to Google Cloud Run or any container platform:

```bash
# Build Docker image
docker build -t hatstack-app .

# Run locally
docker run -p 8080:8080 --env-file .env hatstack-app

# Deploy to Cloud Run
gcloud run deploy hatstack-app --source .
```

The Dockerfile is production-ready and optimized for frontend applications.

## 🆘 Troubleshooting

### "Air not found"

Make sure `$(go env GOPATH)/bin` is in your PATH:

**Linux/Mac:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Windows:**
```cmd
set PATH=%PATH%;%USERPROFILE%\go\bin
```

### "Templ not found"

Same as above - add Go bin directory to PATH.

### "Cannot connect to Datastore"

Check your `.env` file:
- Verify `GOOGLE_APPLICATION_CREDENTIALS` points to valid JSON key
- Or verify `GOOGLE_PROJECT_ID` is correct
- Ensure your service account has Datastore permissions

### Port already in use

Change the port in `.env`:
```env
PORT=3000
```

## 💡 Tips

- Use `./dev.sh` (Linux/Mac) for quick start with auto-setup check
- Keep Air running while developing for instant feedback
- Check `build-errors.log` if builds fail
- Press `Ctrl+C` to stop the dev server

## 🎉 You're Ready!

Visit `http://localhost:8080` and start building your H.A.T. Stack application!

Need help? Check the [full README](../README.md) or open an [issue](https://github.com/runtime-dynamics/hatstack/issues).
