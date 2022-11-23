# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

#COPY *.go ./
COPY ./ ./

RUN go build -o /exchange-message-broker

EXPOSE 7500

CMD [ "/exchange-message-broker" ]