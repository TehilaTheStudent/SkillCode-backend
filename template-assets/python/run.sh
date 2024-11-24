#!/bin/sh

# This script runs the provided Python file using Python3

# Check if a file was provided as an argument
if [ -z "$1" ]; then
  echo "Error: No Python file provided."
  exit 1
fi

FILE="$1"

# Ensure the file exists
if [ ! -f "$FILE" ]; then
  echo "Error: File $FILE does not exist."
  exit 1
fi

# Run the Python file using Python3
python3 "$FILE"
EXIT_CODE=$?

# Exit with the status code from Python3
exit $EXIT_CODE
