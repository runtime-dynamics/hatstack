#!/bin/bash
set -e

echo "==================================="
echo "H.A.T. Stack Development Setup"
echo "==================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"
echo ""

# Check if Air is installed
echo "Checking for Air (live-reload tool)..."
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}Air not found. Installing...${NC}"
    go install github.com/air-verse/air@latest
    
    # Check if installation was successful
    if ! command -v air &> /dev/null; then
        echo -e "${RED}Warning: Air was installed but not found in PATH${NC}"
        echo "Make sure $(go env GOPATH)/bin is in your PATH"
        echo "Add this to your ~/.bashrc or ~/.zshrc:"
        echo "  export PATH=\$PATH:$(go env GOPATH)/bin"
        exit 1
    fi
    echo -e "${GREEN}✓ Air installed successfully${NC}"
else
    echo -e "${GREEN}✓ Air is already installed${NC}"
fi
echo ""

# Check if Templ is installed
echo "Checking for Templ (template generator)..."
if ! command -v templ &> /dev/null; then
    echo -e "${YELLOW}Templ not found. Installing...${NC}"
    go install github.com/a-h/templ/cmd/templ@latest
    
    # Check if installation was successful
    if ! command -v templ &> /dev/null; then
        echo -e "${RED}Warning: Templ was installed but not found in PATH${NC}"
        echo "Make sure $(go env GOPATH)/bin is in your PATH"
        echo "Add this to your ~/.bashrc or ~/.zshrc:"
        echo "  export PATH=\$PATH:$(go env GOPATH)/bin"
        exit 1
    fi
    echo -e "${GREEN}✓ Templ installed successfully${NC}"
else
    echo -e "${GREEN}✓ Templ is already installed${NC}"
fi
echo ""

# Check for .env file
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}.env file not found. Creating template...${NC}"
    cat > .env << 'EOF'
# Google Cloud Configuration
DATASTORE_NAME=your-datastore-name
PUBSUB_TOPIC=your-pubsub-topic
PUBSUB_SUBSCRIPTION=your-pubsub-subscription

# Frontend Configuration
FRONTEND_ENDPOINT=http://local.nitecon.net:8080

# Google Cloud Authentication
# Option 1: Use service account key file
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json

# Option 2: Or use project ID for default credentials
# GOOGLE_PROJECT_ID=your-project-id

# Development Settings
DEBUG=true
PORT=8080
EOF
    echo -e "${GREEN}✓ Created .env file${NC}"
    echo ""
    echo -e "${YELLOW}IMPORTANT: Please edit .env and configure your settings${NC}"
    echo "Required variables:"
    echo "  - DATASTORE_NAME"
    echo "  - PUBSUB_TOPIC"
    echo "  - PUBSUB_SUBSCRIPTION"
    echo "  - FRONTEND_ENDPOINT"
    echo "  - GOOGLE_APPLICATION_CREDENTIALS or GOOGLE_PROJECT_ID"
    echo ""
    echo "Opening .env in default editor..."
    ${EDITOR:-nano} .env
else
    echo -e "${GREEN}✓ .env file exists${NC}"
fi
echo ""

# Setup .air.toml for Unix
echo "Setting up .air.toml for Unix environment..."
if [ -f ".air.unix.toml" ]; then
    cp .air.unix.toml .air.toml
    echo -e "${GREEN}✓ Configured .air.toml for Unix${NC}"
else
    echo -e "${RED}Warning: .air.unix.toml not found${NC}"
fi
echo ""

# Make run script executable
if [ -f "scripts/run-with-env.sh" ]; then
    chmod +x scripts/run-with-env.sh
    echo -e "${GREEN}✓ Made run-with-env.sh executable${NC}"
fi
echo ""

echo "==================================="
echo -e "${GREEN}Setup Complete!${NC}"
echo "==================================="
echo ""
echo "To start the development server, run:"
echo -e "${GREEN}  air${NC}"
echo ""
echo "The server will automatically:"
echo "  1. Generate templ templates"
echo "  2. Load environment variables from .env"
echo "  3. Build and start the application"
echo "  4. Rebuild on file changes"
echo ""
