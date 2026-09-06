# 開発用ステージ（docker-compose等）
FROM golang:1.26.5-alpine AS dev
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
CMD ["go", "run", "."]

# 本番ビルド用ステージ
FROM dev AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# 本番実行用ステージ（Cloud Run等・軽量イメージ）
FROM alpine:3.20 AS prod
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/server .
COPY --from=builder /app/contents ./contents
ENV PORT=8080
EXPOSE 8080
CMD ["./server"]