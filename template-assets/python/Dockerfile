# Use Python 3.10 slim as the base image
FROM python:3.10-slim

# Set the working directory inside the container
WORKDIR /sandbox/app

# Copy the Python template-assets into the container
COPY python/ /sandbox/app/
COPY feedback_schema.json /sandbox/
# Install dependencies from requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Ensure UTF-8 encoding and unbuffered output
ENV PYTHONUNBUFFERED=1 \
    PYTHONIOENCODING=UTF-8
