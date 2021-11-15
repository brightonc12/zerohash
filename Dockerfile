FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build


# Copy and download dependencies
COPY go.mod .
RUN go mod download

# Copy the rest of the code into the container
COPY . .

# Build the application
RUN go build -ldflags="-w -s" -o  dist/main cmd/server.go

FROM alpine:3.13

COPY --from=builder /build/dist/main /main

ENTRYPOINT ["/main"]
