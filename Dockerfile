FROM golang:1.22.5-alpine as builder
WORKDIR /builder
ENV GO111MODULE=on CGO_ENABLED=0

COPY . .
RUN go build -ldflags "-extldflags '-static' -s -w" -o /builder/main /builder/src/main.go


FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder /builder/.env.example /app/.env
COPY --from=builder /builder/main main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
