FROM python:3.10-slim

# Set a non-root user for added security
RUN useradd -ms /bin/bash sandbox-user

# Create a directory for the user code
WORKDIR /sandbox

# Switch to the sandbox user
USER sandbox-user

# Copy the Python assets into the container
COPY assets/python/run_tests.py /sandbox/run_tests.py
COPY assets/python/evaluator.py /sandbox/evaluator.py
COPY assets/python/ds_utils.py /sandbox/ds_utils.py

# Ensure UTF-8 encoding and unbuffered output
ENV PYTHONUNBUFFERED=1 \
    PYTHONIOENCODING=UTF-8

# Command to run the tests
CMD ["python3", "/sandbox/run_tests.py"]
