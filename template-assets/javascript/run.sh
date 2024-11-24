#!/bin/sh

# This script runs the provided JavaScript file using Node.js

# Check if a file was provided as an argument
if [ -z "$1" ]; then
  echo "Error: No JavaScript file provided."
  exit 1
fi

FILE="$1"

# Ensure the file exists
if [ ! -f "$FILE" ]; then
  echo "Error: File $FILE does not exist."
  exit 1
fi

# Run the JavaScript file using Node.js
node "$FILE"
EXIT_CODE=$?

# Exit with the status code from Node.js
exit $EXIT_CODE
