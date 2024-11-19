FROM golang:1.22-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o main .

EXPOSE 8080

ENV PREFIX_SERVICE=SERVER_ADDRESS_

ENV SERVER_ADDRESS_1=http://host.docker.internal:3000/server

ENV SERVER_ADDRESS_2=http://host.docker.internal:3001

CMD ["./main"]

 