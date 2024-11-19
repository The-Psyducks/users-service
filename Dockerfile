# Build stage
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /home/app

COPY /server/go.mod /server/go.sum ./
RUN go mod download

COPY /server ./

RUN go build -o twitsnap ./main.go

# Test stage
FROM builder AS twitsnap-test-stage

# CMD ["go", "test", "-v", "./tests"]
COPY /server/rsc/twitsnap.png ./tests/rsc/twitsnap.png
CMD ["go", "test", "-cover", "-coverprofile=coverage/coverage.out", ".", "./...", "-v"]

# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/twitsnap ./
COPY /server/rsc/twitsnap.png ./rsc/twitsnap.png

# Create a service account file if it doesn't exist
RUN test -f /home/app/service-account.json || echo '{}' > /home/app/service-account.json

ENTRYPOINT ["./twitsnap"]
