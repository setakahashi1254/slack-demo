# --- ステージ1: ビルド環境 ---
FROM golang:1.26-alpine AS builder
WORKDIR /app

# 依存関係のキャッシュとインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてコンパイル
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- ステージ2: 実行環境 ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# ビルドされたバイナリだけをコピー
COPY --from=builder /app/main .

# アプリケーションの実行
CMD ["./main"]