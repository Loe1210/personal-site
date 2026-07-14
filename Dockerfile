# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder
WORKDIR /app

ARG SERVICE_PATH=./
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn
ARG GOPRIVATE=gitlab.tanwan.com

ENV GOPROXY=$GOPROXY
ENV GOSUMDB=$GOSUMDB
ENV GOPRIVATE=$GOPRIVATE

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ${SERVICE_PATH}

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/app /app/app
COPY configs /app/configs
COPY services /app/services
COPY static /app/static

EXPOSE 8888 9001 9002 9003 9004

CMD ["/app/app"]