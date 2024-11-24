FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /home/app

COPY /server/go.mod /server/go.sum ./
RUN go mod download

COPY /server ./

RUN go build -o twitsnap ./main.go

# Test stage
FROM builder AS twitsnap-test-stage

# Crear carpeta para el reporte de cobertura
RUN mkdir -p /home/app/coverage

# Establecer permisos de escritura/lectura para la carpeta coverage
RUN chmod -R 777 /home/app/coverage

# Usar sh -c para ejecutar el comando completo de pruebas y generación del reporte HTML
# CMD ["sh", "-c", "go test -cover -v -coverprofile=/home/app/coverage/coverage.out ./... && go tool cover -html=/home/app/coverage/coverage.out -o /home/app/coverage/coverage.html && chmod -R 777 /home/app/coverage"]
CMD ["sh", "-c", "go test -cover -coverprofile=/home/app/coverage/coverage.out ./... && go tool cover -html=/home/app/coverage/coverage.out -o /home/app/coverage/coverage.html && chmod -R 777 /home/app/coverage"]