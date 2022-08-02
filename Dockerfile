FROM golang:latest as builder

RUN mkdir -p /auth
ADD . /auth
WORKDIR /auth

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /auth ./cmd/main.go

FROM scratch
COPY --from=builder /auth /auth
COPY --from=builder /etc/ssl/certs /etc/ssl/certs/
WORKDIR /auth

CMD ["./main"]
EXPOSE 3000 4000
