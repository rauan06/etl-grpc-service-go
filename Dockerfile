# syntax=docker/dockerfile:1

FROM golang:1.24.2

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project (so we can build from cmd/etl)
COPY . .

# Build the ETL binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /etl-service ./cmd/etl

# Expose the port your service listens on
EXPOSE 8099

# Run the binary
CMD ["/etl-service"]
