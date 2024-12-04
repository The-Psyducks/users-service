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
CMD ["sh", "-c", "go test -cover -coverprofile=coverage/coverage.tmp.out $(go list ./... | grep -v 'src/auth') && \
cat coverage/coverage.tmp.out | grep -v -e 'main.go' -e 'rabbit_messages.go' -e 'send_email.go' > coverage/coverage.out"]


# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/twitsnap ./
COPY /server/rsc/twitsnap.png ./rsc/twitsnap.png

# Create a service account file if it doesn't exist
RUN test -f /home/app/service-account.json || echo '{}' > /home/app/service-account.json

ENTRYPOINT ["./twitsnap"]