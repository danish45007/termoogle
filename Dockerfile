FROM golang:latest

LABEL maintainer="danish <danish45007@gmail.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV API_KEY_KUTT $API_KEY_KUTT

RUN go build

CMD [ "./my-google.com" ]