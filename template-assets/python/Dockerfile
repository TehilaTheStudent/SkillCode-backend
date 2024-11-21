# Use Python 3.10 slim as the base image
FROM python:3.10-slim

# Set the working directory inside the container
WORKDIR /sandbox

# Copy the Python template-assets into the container
COPY . /sandbox/

# Install dependencies from requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Ensure UTF-8 encoding and unbuffered output
ENV PYTHONUNBUFFERED=1 \
    PYTHONIOENCODING=UTF-8


