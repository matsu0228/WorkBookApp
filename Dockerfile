FROM golang:1.14
#FROM golang:1.14-alpine
# コンテナログイン時のディレクトリ指定
WORKDIR /opt/sandbox-docker-compose-go
# ホストのファイルをコンテナの作業ディレクトリにコピー
COPY . .
# ADD . .
# ビルド
RUN go build -o app main.go
# 起動
CMD ["/opt/sandbox-docker-compose-go/app"]
