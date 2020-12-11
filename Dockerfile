FROM golang:1.15-alpine AS builder
WORKDIR /app
COPY . .
RUN go get ./
RUN go build -o server main.go

FROM alpine
COPY --from=builder /app/server .
ENTRYPOINT [ "/server" ]