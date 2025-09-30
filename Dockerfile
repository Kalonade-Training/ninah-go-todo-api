FROM golang:1.22-alpine AS build

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy everything
COPY . .

# Download and build
RUN go mod download
RUN go build -o main .

EXPOSE 8091

CMD ["./main"]