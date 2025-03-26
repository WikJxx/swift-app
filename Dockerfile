FROM golang:1.24.1-alpine

WORKDIR /app

COPY app/go.mod app/go.sum ./

RUN go mod tidy

COPY app/ .

RUN go build -o swift-app .

EXPOSE 8080

CMD ["./swift-app"]
