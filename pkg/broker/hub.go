package broker

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	//message chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	message chan Message

	Subscribe   chan string
	Unsubscribe chan string
}

type ClientMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
	/*"method": "SUBSCRIBE",
	"params": [
		"btcusdt@aggTrade",
		"btcusdt@depth"
	],
	"id": 1*/
}

type Message struct {
	Client        *Client
	ClientMessage ClientMessage
}

type ResponseMessage struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		message:     make(chan Message),
		Subscribe:   make(chan string),
		Unsubscribe: make(chan string),
	}
}

func (h *Hub) PublishMessage(topic string, message interface{}) {
	logrus.Tracef("Publish message on topic %s", topic)

	eb.Publish(topic, message)
}

func (h *Hub) processMessage(message Message) {
	switch strings.ToUpper(message.ClientMessage.Method) {
	case "SUBSCRIBE":
		h.processSubscribeMessage(message.Client, message.ClientMessage)

	case "UNSUBSCRIBE":
		h.processUnSubscribeMessage(message.Client, message.ClientMessage)

	case "PUBLISH":
		processPublishMessage(message.Client, message.ClientMessage)

	default:
		logrus.Trace("Unknown message")
	}
}

func processPublishMessage(client *Client, message ClientMessage) {
	logrus.Tracef("Publish message  from client %d on topic %s: %s", client.id, message.Params[0], message.Params[1])
	eb.Publish(message.Params[0], message.Params[1])
}

func (h *Hub) processSubscribeMessage(client *Client, message ClientMessage) {
	logrus.Info("SUBSCRIBE", message)
	for _, topic := range message.Params {
		_, found := client.getChannel(topic)
		if found {
			logrus.Tracef("Client %d already subscribed on %s\n", client.id, topic)
			continue
		}
		logrus.Tracef("Client %d subscribing on %s\n", client.id, topic)
		//bus.Subscribe(topic, calculator)
		ch := make(chan DataEvent)
		go client.subClient(topic, ch)

		eb.Subscribe(topic, ch)
		h.Subscribe <- topic
	}
	client.sendMessage(ResponseMessage{
		Id:     message.Id,
		Result: "success",
	})
}

func (h *Hub) processUnSubscribeMessage(client *Client, message ClientMessage) {
	logrus.Trace("UNSUBSCRIBE", message)

	for _, v := range message.Params {
		h.unSubscribe(v, client)
	}
	client.sendMessage(ResponseMessage{
		Id:     message.Id,
		Result: "success",
	})
}

func (h *Hub) unSubscribe(topic string, client *Client) {
	ch, _ := client.getChannel(topic)
	if ch != nil {
		eb.Unsubscribe(topic, ch)
		client.unsubClient(topic)
		logrus.Tracef("Client %d on topic %s unsubscribed", client.id, topic)
		h.Unsubscribe <- topic
	} else {
		logrus.Tracef("There is no channel on top %s on client %d", topic, client.id)
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			logrus.Trace("Close in broadcastMessage")
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.processRegister(client)

		case client := <-h.unregister:
			h.processUnregister(client)

		case message := <-h.message:
			h.processMessage(message)
		}
	}
}

func (h *Hub) processRegister(client *Client) {
	logrus.Tracef("client %d registered", client.id)
	h.clients[client] = true
	h.printClients()
}

func (h *Hub) processUnregister(client *Client) {
	logrus.Tracef("client %d unregistered", client.id)
	if _, ok := h.clients[client]; ok {
		for topic := range client.channels {
			h.unSubscribe(topic, client)
		}
		delete(h.clients, client)
		close(client.send)
		h.printClients()
	}
}

func (h *Hub) printClients() {
	logrus.Trace("Number of clients: ", len(h.clients))
}

var clientId = 0

// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	clientId++
	client := &Client{id: clientId, hub: h, conn: conn, send: make(chan []byte, 256), channels: make(map[string]chan DataEvent)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
