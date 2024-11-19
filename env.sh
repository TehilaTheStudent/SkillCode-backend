#!/bin/bash

# Dynamically set PROJECT_ROOT to the current directory,
# This is for local development
export PROJECT_ROOT=$(pwd)
export MODE_ENV=development

echo "PROJECT_ROOT set to $PROJECT_ROOT"
echo "MODE_ENV set to $MODE_ENV"
# Usage:
# source env.sh
