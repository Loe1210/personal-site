# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/personal-site ./

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/personal-site /app/personal-site
COPY static /app/static
COPY configs /app/configs

EXPOSE 8888

CMD ["/app/personal-site", "-config", "/app/configs/config.yaml"]
