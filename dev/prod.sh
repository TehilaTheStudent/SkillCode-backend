#!/bin/bash


export MODE_ENV=production
# Dynamically set PROJECT_ROOT to the current directory
export PROJECT_ROOT=$(pwd)

echo "MODE_ENV set to $MODE_ENV"
