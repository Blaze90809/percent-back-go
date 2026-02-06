# Stage 1: Build the React frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built frontend assets from Stage 1 into the backend's static directory
COPY --from=frontend-builder /app/frontend/build ./static
RUN go build -o /react-app-golang

# Stage 3: Final lightweight image
FROM alpine:latest
WORKDIR /root/
# Copy the binary and the static files from the backend builder
COPY --from=backend-builder /react-app-golang .
COPY --from=backend-builder /app/static ./static
# Ensure .env is available if needed, or use environment variables in production
# COPY --from=backend-builder /app/.env .

EXPOSE 8080
CMD ["sh", "-c", "./react-app-golang"]
