# Używamy obrazu Go 1.24.1 jako bazowego
FROM golang:1.24.1-alpine

# Ustawiamy katalog roboczy
WORKDIR /app

# Kopiujemy pliki go.mod i go.sum, aby zainstalować zależności
COPY go.mod go.sum ./

# Instalujemy zależności
RUN go mod tidy

# Kopiujemy cały kod źródłowy do katalogu roboczego
COPY . .

# Budujemy aplikację
RUN go build -o swift-app .

# Eksponujemy port, na którym aplikacja będzie działać
EXPOSE 8080

# Uruchamiamy aplikację
CMD ["./swift-app"]
