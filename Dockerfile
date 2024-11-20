FROM golang:1.22-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o main .

EXPOSE 8080

ENV PREFIX_SERVICE=SERVER_ADDRESS_

ENV SERVER_ADDRESS_1=http://192.168.13.225:6969

ENV SERVER_ADDRESS_2=http://192.168.13.147:6969

ENV SERVER_ADDRESS_3=http://host.docker.internal:3000

CMD ["./main"]

 