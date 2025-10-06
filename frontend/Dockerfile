# Multi-stage Dockerfile for Vue frontend (Vite build + Nginx serve)

# 1) Builder stage
FROM node:20-alpine AS builder
WORKDIR /app

# Install dependencies
COPY frontend/package*.json ./
RUN npm ci || npm install

# Copy sources and build
COPY frontend/ ./
RUN npm run build

# 2) Runtime stage (Nginx)
FROM nginx:1.25-alpine

# Remove default config and add our SPA + API proxy config
RUN rm -f /etc/nginx/conf.d/default.conf
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
