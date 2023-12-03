# syntax=docker/dockerfile:1

FROM golang:1.21.4 AS builder

ENV PORT_API=8080
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source into container
COPY . .

# Build
RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -o ingestor ./cmd/ingestor/main.go

# Host container
FROM ubuntu:22.04
EXPOSE $PORT_API
COPY --from=builder /app/ingestor /usr/local/bin/
CMD ["ingestor"]