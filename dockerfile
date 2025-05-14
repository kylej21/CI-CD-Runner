FROM golang:1.24.3

RUN apt-get update && \
    apt-get install -y curl git && \
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

WORKDIR /app
COPY . .

RUN go build -o api ./cmd/apiserver
EXPOSE 8080
ENTRYPOINT ["./api"]
