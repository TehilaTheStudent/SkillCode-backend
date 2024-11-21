# Use Node.js 18 Alpine as the base image
FROM node:18-alpine


# Set the working directory
WORKDIR /sandbox/app

# Copy JavaScript template-assets to the container
COPY javascript/ /sandbox/app/
COPY feedback_schema.json /sandbox/

RUN npm install --production


# Set environment variables
ENV NODE_ENV=production
