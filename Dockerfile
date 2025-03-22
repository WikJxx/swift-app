FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o swift-app .

EXPOSE 8080

CMD ["./swift-app"]
