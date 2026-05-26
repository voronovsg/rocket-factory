# Build context must be the repository root (paths below are relative to it):
#   docker build -f deploy/docker/inventory/Dockerfile .
FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./

COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY assembly/go.mod assembly/go.sum ./assembly/
COPY notification/go.mod notification/go.sum ./notification/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/

RUN go mod download

COPY platform ./platform
COPY shared ./shared
COPY order ./order
COPY inventory ./inventory
COPY payment ./payment
COPY assembly ./assembly
COPY notification ./notification

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.49/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/main.go

FROM alpine:3.22.4

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-inventory .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

COPY --from=builder /app/shared/api/inventory/v1 /app/shared/api/inventory/v1

ENV INVENTORY_HTTP_STATIC_DIR=/app/shared/api/inventory/v1

EXPOSE 50051

ENTRYPOINT ["./app-inventory"]
