FROM golang:1.26.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o minimart-api .

EXPOSE 8080

CMD ["./minimart-api"]