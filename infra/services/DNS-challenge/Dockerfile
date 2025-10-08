FROM golang:1.25.1-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0

COPY . .

RUN go mod tidy

RUN go build -o /beget-hook

FROM alpine:3.22.1

RUN apk add --no-cache ca-certificates

COPY --from=build /beget-hook /usr/local/bin/beget-hook

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/beget-hook"]
