FROM golang:1.25-alpine

WORKDIR /app

# Copy only go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the API binary
RUN go build -o api ./cmd/api

CMD ["./api"]
