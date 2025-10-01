FROM golang:1.23-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy everything
COPY . .

# Download and build
RUN go mod download
RUN go build -o main ./cmd/main.go

EXPOSE 8091

CMD ["./main"]