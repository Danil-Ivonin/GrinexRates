FROM golang:1.25.5 AS builder

WORKDIR /app

# Copy dependency files first for layer-cache optimization
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source tree
COPY . .

# Build a fully static, stripped binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd

# Stage 2: fminimal distroless image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /home/nonroot

COPY --from=builder /server /server
COPY config.yml .

# gRPC port
EXPOSE 50051
# Prometheus metrics port
EXPOSE 9090

ENTRYPOINT ["/server"]