# Go-ney Dockerfile
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copiar go mod y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagen final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario desde builder
COPY --from=builder /app/main .

# Crear usuario no-root
RUN adduser -D -s /bin/sh goney
USER goney

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["./main"]
