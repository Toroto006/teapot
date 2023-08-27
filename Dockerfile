# First build the go executable
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o teapot .

# Second build the image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/teapot .

EXPOSE 8080

CMD ["./teapot"]