FROM golang:1.23 AS builder

RUN mkdir -p /go/app

WORKDIR /go/app

COPY go.mod go.mod
COPY go.sum go.sum
COPY server/ server/
COPY api/go/ api/go/

WORKDIR server

RUN --mount=type=ssh CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main

FROM alpine

COPY --from=builder /main ./

CMD ["./main"]
