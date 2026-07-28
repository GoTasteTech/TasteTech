# ===== STAGE 1: Build =====
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o tastetech .

# ===== STAGE 2: Run (image kecil) =====
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage build
COPY --from=builder /app/tastetech .

# Copy folder view, static (template HTML & aset)
COPY --from=builder /app/view ./view
COPY --from=builder /app/static ./static

# Port default
EXPOSE 8000

# Jalankan app
CMD ["./tastetech"]
