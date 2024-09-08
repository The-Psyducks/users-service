# Build stage
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /home/app

COPY /server/go.mod /server/go.sum ./
RUN go mod download

COPY /server ./

RUN go build -o twitsnap ./src/main.go

# Test stage
FROM builder AS twitsnap-test-stage

CMD ["go", "test", "-v", "./tests"]

# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/twitsnap ./

CMD ["./twitsnap"]