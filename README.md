# Introduction
Real-time message broker for Binance.
Exchange Message Broker is a Golang application that serves WebSocket connections. This application implements the Pub/Sub 
pattern for making high-performance concurrent connections. Clients connected to the server
and subscribe to topics, and the server creates a subscription on Binance. if nobody subscribed to a specific topic, the server
unsubscribe from that topic on Binance.

Also, there is an HTML file that has Javascript client implementation that connects to this server and subscribes to some topics.

## Run
You should edit appsettings.json file and set remote server port for listening
Run this command to get all necessary packages:
### `go get all`

## Test
Use this command for unit testing:
### `go test ./...`