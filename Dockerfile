FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o code-processor .

CMD ["./code-processor"]

EXPOSE 8080
