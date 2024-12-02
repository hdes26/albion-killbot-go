# Etapa 1: Construcción de la aplicación
FROM golang:1.23.2-alpine as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el código fuente a la carpeta de trabajo
COPY . .

# Construir el binario de la aplicación
RUN go build -o main cmd/app/main.go

# Etapa 2: Crear la imagen de producción final
FROM alpine:latest

# Instalar ca-certificates para las conexiones HTTPS
RUN apk --no-cache add ca-certificates

# Establecer el directorio de trabajo en el contenedor final
WORKDIR /root/

# Copiar el binario de la etapa de construcción
COPY --from=builder /app/main .

# Comando para ejecutar el binario
CMD ["./main"]