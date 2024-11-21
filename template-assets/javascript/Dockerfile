# Use Node.js 18 Alpine as the base image
FROM node:18-alpine


# Set the working directory
WORKDIR /sandbox

# Copy JavaScript template-assets to the container
COPY . /sandbox/

RUN npm install --production


# Set environment variables
ENV NODE_ENV=production
