#!/bin/bash

# Quick development server launcher
# Ensures setup is complete before running air

if [ ! -f ".env" ]; then
    echo "No .env file found. Running setup..."
    ./setup.sh
fi

if [ ! -f ".air.toml" ]; then
    echo "No .air.toml found. Running setup..."
    ./setup.sh
fi

# Make sure run script is executable
if [ -f "scripts/run-with-env.sh" ]; then
    chmod +x scripts/run-with-env.sh
fi

# Run air
air
