# BUILD EXECUTABLE
FROM golang:1.22.1-alpine3.19 as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

RUN chmod +x /app/authApp

# DEPLOY
FROM alpine:3.14

RUN mkdir /app

COPY --from=builder /app/authApp /app

CMD ["/app/authApp"]
