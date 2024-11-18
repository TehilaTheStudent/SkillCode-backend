# Use Node.js as the base image
FROM node:18-alpine

# Set a non-root user for added security
RUN useradd -ms /bin/bash sandbox-user

# Create a working directory
WORKDIR /sandbox

# Switch to the sandbox user
USER sandbox-user

# Copy the JavaScript assets (evaluator, utils, tests)
COPY assets/javascript/evaluator.js /sandbox/evaluator.js
COPY assets/javascript/ds_utils.js /sandbox/ds_utils.js
COPY assets/javascript/run_tests.js /sandbox/run_tests.js

# Set environment variables
ENV NODE_ENV=production

# Command to run the tests
CMD ["node", "/sandbox/run_tests.js"]
