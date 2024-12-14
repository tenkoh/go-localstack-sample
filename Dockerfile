# 説明の簡単化のためマルチステージビルドは使わない
FROM golang:1.23-bullseye

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

CMD ["./main"]
