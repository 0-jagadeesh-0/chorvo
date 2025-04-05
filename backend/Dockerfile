# base image
FROM golang:1.23-alpine

WORKDIR /backend

COPY go.* ./

RUN go mod download

COPY . .

# Copy the .env file to the working directory
COPY .env .env

RUN go build -o backend ./cmd/main.go

EXPOSE 8080

CMD ["./backend"]