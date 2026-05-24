FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./

COPY . .

RUN go mod tidy

RUN go build -o notif_app ./notification_service/main.go

EXPOSE 8081

CMD ["./notif_app"]