# Build context must be the repository root:
#   docker build -f platform/pkg/testcontainers/stubiam/Dockerfile .
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
COPY iam/go.mod iam/go.sum ./iam/

RUN go mod download

COPY platform ./platform
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stub-iam ./platform/pkg/testcontainers/stubiam/cmd/main.go

FROM alpine:3.22.4

WORKDIR /app

COPY --from=builder /app/stub-iam .

ENV GRPC_HOST=0.0.0.0
ENV GRPC_PORT=50053

EXPOSE 50053

ENTRYPOINT ["./stub-iam"]
