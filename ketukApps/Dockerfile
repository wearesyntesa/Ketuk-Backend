# Stage 1: Build stage
FROM golang:1.25.1-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy go mod dan go sum files untuk install dependencies
# Ini dilakukan sebelum copy source code agar caching docker layer optimal
COPY go.mod go.sum ./

# Download semua dependencies
RUN go mod download

# Copy source code dari folder ketukApps ke dalam container
COPY . .

# Build aplikasi Go menjadi binary bernama 'main'
# CGO_ENABLED=0 diperlukan agar binary statis dan bisa jalan di alpine murni
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Production stage (Image akhir yang akan dideploy)
FROM alpine:latest

# Install ca-certificates jika aplikasi Anda melakukan request HTTPS ke luar
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary hasil build dari stage builder
COPY --from=builder /app/main .

# Expose port (sesuaikan dengan port aplikasi Anda, misal 8080)
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["./main"]
