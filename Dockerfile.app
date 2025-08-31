# Created by an LLM.
# Edited by me: docker image versions.

# Frontend React/Vite application Dockerfile
# Based on the Justfile dev_app configuration

# Build stage
FROM node:24-alpine AS builder

WORKDIR /build

# Copy package files for dependency caching
COPY app/package.json app/package-lock.json ./
RUN npm ci --only=production

# Copy source code
COPY app/src ./src/
COPY app/index.html ./
COPY app/tsconfig*.json ./
COPY app/vite.config.ts ./
COPY app/config.js ./
COPY app/config.d.ts ./

# Build the application
RUN npm run build

# Production stage
FROM nginx:1.29.1-alpine

# Copy built assets
COPY --from=builder /build/dist /usr/share/nginx/html

# Copy nginx configuration if needed
# COPY app/nginx.conf /etc/nginx/nginx.conf

# Create a simple nginx config for SPA
RUN echo 'server { \
    listen 8080; \
    server_name localhost; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
    location /api/ { \
        proxy_pass http://api:8081/; \
        proxy_set_header Host $host; \
        proxy_set_header X-Real-IP $remote_addr; \
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; \
        proxy_set_header X-Forwarded-Proto $scheme; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --start-interval=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
