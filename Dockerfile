# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser to run the application
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -tags netgo -installsuffix netgo \
    -o /app/fast-cc-hooks \
    ./cmd/fast-cc-hooks

# Final stage - minimal image
FROM scratch

# Copy certificates, timezone data, and user info from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary
COPY --from=builder /app/fast-cc-hooks /fast-cc-hooks

# Use non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["/fast-cc-hooks"]
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="fast-cc-hooks"
LABEL org.opencontainers.image.description="Fast conventional commits git hooks"
LABEL org.opencontainers.image.source="https://github.com/stevengreensill/fast-cc-git-hooks"
LABEL org.opencontainers.image.licenses="MIT"