# -------------------------------------------------------------
# 開発用ステージ（docker-compose等）
# -------------------------------------------------------------
FROM golang:1.26.5-alpine AS dev
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
CMD ["go", "run", "."]

# -------------------------------------------------------------
# 本番ビルド用ステージ（Goバイナリのビルド）
# -------------------------------------------------------------
FROM dev AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# -------------------------------------------------------------
# コンテンツ取得ステージ（別リポジトリから clone / ローカルフォールバック）
# -------------------------------------------------------------
FROM alpine:3.20 AS content-fetcher
RUN apk add --no-cache git

# ビルド引数
ARG DOCS_REPO_URL=""
ARG DOCS_BRANCH="main"

WORKDIR /app
# ローカルの contents をコピー（DOCS_REPO_URL 未指定時のフォールバック用）
COPY contents ./contents-local
RUN if [ -n "$DOCS_REPO_URL" ]; then \
      git clone --depth 1 -b ${DOCS_BRANCH} ${DOCS_REPO_URL} ./contents; \
    else \
      cp -r ./contents-local ./contents; \
    fi

# -------------------------------------------------------------
# 本番実行用ステージ（Cloud Run 等・軽量イメージ）
# -------------------------------------------------------------
FROM alpine:3.20 AS prod
WORKDIR /app

# ユーザーを作成（セキュリティ強化）
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -g '' appuser

# ビルドしたバイナリと contents をコピー
COPY --from=builder /app/server .
COPY --from=content-fetcher /app/contents ./contents

# 実行ユーザーに権限を付与
RUN chown -R appuser:appuser /app

USER appuser
ENV PORT=8080
EXPOSE 8080
CMD ["./server"]