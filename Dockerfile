# Step 1: Build the Go binary
FROM golang:1.22 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Step 2: Run with a minimal image
FROM gcr.io/distroless/base-debian12

COPY --from=build /app/server /server

# Cloud Run expects your app to listen on PORT env var
ENV PORT=8080
CMD ["/server"]
