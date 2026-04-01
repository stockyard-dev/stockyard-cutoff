FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go mod download && CGO_ENABLED=0 go build -o cutoff ./cmd/cutoff/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/cutoff .
ENV PORT=9330 DATA_DIR=/data
EXPOSE 9330
CMD ["./cutoff"]
