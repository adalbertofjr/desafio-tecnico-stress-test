# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copiar código fonte
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o stresstest ./cmd/stresstest

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário do builder
COPY --from=builder /app/stresstest .

ENTRYPOINT ["./stresstest"]
