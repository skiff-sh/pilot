FROM golang:1.23 AS builder

ARG GRPC_HEALTH_CHECK_VERSION="v0.4.36"

ARG TARGETARCH

RUN wget https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_CHECK_VERSION}/grpc_health_probe-linux-${TARGETARCH} && \
    mv grpc_health_probe-linux-${TARGETARCH} grpc_health_probe && \
    chmod +x grpc_health_probe

RUN mkdir -p /go/app
WORKDIR /go/app

COPY go.mod go.mod
COPY go.sum go.sum
COPY server/ server/
COPY pkg/ pkg/
COPY api/go/ api/go/

WORKDIR server

RUN --mount=type=ssh CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main

FROM alpine

COPY --from=builder /main ./

COPY --from=builder /go/grpc_health_probe /grpc_health_probe

CMD ["./main"]
