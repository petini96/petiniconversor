FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o youtube_downloader.exe main.go

FROM alpine:latest as final

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/youtube_downloader.exe /youtube_downloader.exe

CMD ["echo", "Este container não executa nada, apenas cria o binário"]
