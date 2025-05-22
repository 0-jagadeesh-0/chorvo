# base image
FROM golang:1.23-alpine

WORKDIR /server

COPY go.* ./

RUN go mod download

COPY . .

# Copy the .env file to the working directory
COPY .env .env

RUN go build -o server ./cmd/server/main.go

EXPOSE 8080

CMD ["./server"]