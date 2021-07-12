package bridge

import (
	"github.com/sirupsen/logrus"
	"talox.ai/exchangeMessageBroker/pkg/broker"
	"talox.ai/exchangeMessageBroker/pkg/exchange"
)

type Bridge struct {
	hub         *broker.Hub
	exchange    *exchange.Exchange
	subscribers map[string]int
}

func NewBridge(
	hub *broker.Hub,
	exchange *exchange.Exchange,
) *Bridge {
	return &Bridge{
		hub:         hub,
		exchange:    exchange,
		subscribers: make(map[string]int),
	}
}

func (b *Bridge) messageHandler(topic string, data string) {
	b.hub.PublishMessage(topic, data)
}

func (b *Bridge) Run() {
	for {
		select {
		case sTopic := <-b.hub.Subscribe:
			logrus.Tracef("Bridge Subscribe on topic %s", sTopic)
			count, ok := b.subscribers[sTopic]
			if count == 0 {
				_, err := b.exchange.Subscribe(sTopic, b.messageHandler)
				if err != nil {
					logrus.Errorf("Binance, error in subscription on topic %s: %s", sTopic, err.Error())
				} else {
					b.subscribers[sTopic] = 1
				}
				continue
			}
			if ok && count > 0 {
				b.subscribers[sTopic]++
				continue
			} else {
				panic("It is not ok and count <= 0")
			}
		case uTopic := <-b.hub.Unsubscribe:
			count, _ := b.subscribers[uTopic]
			if count == 0 {
				continue
			}
			if count == 1 {
				b.subscribers[uTopic] = 0
				b.exchange.UnSubscribe(uTopic)
				continue
			}
			if count > 1 {
				b.subscribers[uTopic]--
				continue
			} else {
				panic("How is this possible!")
			}
		}
	}
	logrus.Info("End of bridge run func")
}
