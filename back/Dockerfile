
FROM golang:1.26-alpine AS builder

WORKDIR /app

# зависимости (УЖЕ В back/)
COPY go.mod go.sum ./
RUN go mod download

# исходники (УЖЕ в back/)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ЭТАП 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV PORT=7860
EXPOSE 7860

CMD ["./main"]
