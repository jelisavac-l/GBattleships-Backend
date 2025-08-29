# Step 1: Build the Go binary
FROM golang:1.24.3 AS build

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (better build cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the server from cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Step 2: Run the binary in a minimal image
FROM gcr.io/distroless/base-debian12

# Copy built binary from builder stage
COPY --from=build /app/server /server

# Railway (and most cloud providers) set PORT dynamically
ENV PORT=8080
EXPOSE 8080

# Run the server
CMD ["/server"]
