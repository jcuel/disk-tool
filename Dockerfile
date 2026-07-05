# Build frontend
FROM node:22-alpine AS web
WORKDIR /app
COPY web/package.json web/package-lock.json* ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Build Go binary
FROM golang:1.25.11-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /app/dist ./cmd/disk-tool/static
RUN CGO_ENABLED=0 go build -o /disk-tool ./cmd/disk-tool

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /disk-tool /usr/local/bin/disk-tool
COPY scripts/smoke-in-container.sh /scripts/smoke-in-container.sh
RUN chmod +x /scripts/smoke-in-container.sh
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["disk-tool"]
CMD ["serve", "--port", "8080", "--no-open"]
