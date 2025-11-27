# Development Scripts

This directory contains scripts for local development. These scripts are excluded from Docker builds via `.dockerignore`.

**📖 Main Documentation:** See the [main README](../README.md) for project overview and getting started.

**🐳 Deployment:** This project includes a production-ready `Dockerfile` that can be deployed directly to Google Cloud Run or any container platform. The development scripts are for local development only and are not included in the Docker build.

## Files

### Windows Scripts

#### setup-air.ps1
PowerShell script that performs one-time setup:
- Creates a `.env` file if it doesn't exist (with default template)
- Checks if Air (live-reload tool) is installed
- Installs Air if not present using `go install`
- Checks if Templ (template generator) is installed
- Installs Templ if not present using `go install`
- Configures `.air.toml` for Windows

Called by `setup.bat` in the project root.

#### run-with-env.ps1
Wrapper script used by Air to:
- Load environment variables from `.env` file
- Execute the compiled application with those variables
- Called automatically on each rebuild

### Linux/Mac Scripts

#### setup-air.sh
Bash script that performs one-time setup:
- Creates a `.env` file if it doesn't exist (with default template)
- Checks if Air (live-reload tool) is installed
- Installs Air if not present using `go install`
- Checks if Templ (template generator) is installed
- Installs Templ if not present using `go install`
- Configures `.air.toml` for Unix systems
- Makes run-with-env.sh executable

Called by `setup.sh` in the project root.

#### run-with-env.sh
Wrapper script used by Air to:
- Load environment variables from `.env` file
- Execute the compiled application with those variables
- Called automatically on each rebuild

## Usage

### First-Time Setup

**Windows:**
```cmd
setup.bat
```

**macOS/Linux:**
```bash
./setup.sh
```

This will:
1. Install Air if not present
2. Install Templ if not present
3. Create `.env` file if it doesn't exist (prompts you to configure it)

### Running the Development Server

After setup, start the development server with:

```bash
air
```

The server will automatically:
1. Generate templ templates (`templ generate`)
2. Load environment variables from `.env`
3. Build the Go application
4. Start the server
5. Rebuild and restart when you make changes to `.go`, `.templ`, or `.html` files

**Automatic Templ Generation:**
Air is configured to run `templ generate` before each build, so you don't need to manually generate templ files. Any changes to `.templ` files will trigger:
1. Templ code generation
2. Go compilation
3. Application restart

## Environment Configuration

On first run, the script will create a `.env` file with default values. You'll be prompted to edit it with your actual configuration. Required variables include:

- `DATASTORE_NAME`
- `PUBSUB_TOPIC`
- `PUBSUB_SUBSCRIPTION`
- `FRONTEND_ENDPOINT`
- `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_PROJECT_ID`

The `.env` file is gitignored and will not be committed to the repository.

## Requirements

**All Platforms:**
- Go must be installed and in your PATH

**Windows:**
- PowerShell execution policy must allow script execution (the batch file handles this)

**Linux/Mac:**
- Bash shell
- Make scripts executable: `chmod +x setup.sh scripts/*.sh`

## Configuration

Air configuration is defined in `.air.toml` in the project root. Platform-specific templates are available:
- `.air.unix.toml` - Template for Linux/Mac (copied to `.air.toml` by setup.sh)
- Windows uses the default `.air.toml` (configured by setup-air.ps1)

The `.air.toml` file is gitignored to allow platform-specific configurations.
