package exchange

import (
	"encoding/json"
	"errors"
	bin "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"hamednikzad/exchange-message-broker/pkg/exchange/binance"
	"strings"
)

type MsgHandler func(topic string, event string)
type ErrHandler func(err error)

// Subscribe topic => "exchange_operation_argument"
// "binance_kline_ltcbtc:1m"
// "binance_depth_ltcbtc"
type Exchange struct {
	subscriptions map[string]chan struct{}
}

func NewExchange() *Exchange {
	return &Exchange{
		subscriptions: make(map[string]chan struct{}),
	}
}

func (e *Exchange) UnSubscribe(topic string) {
	if c, ok := e.subscriptions[topic]; ok {
		logrus.Infof("Subscription on %s exists and closed", topic)
		delete(e.subscriptions, topic)
		c <- struct{}{}
		close(c)
	}
}

func (e *Exchange) Subscribe(topic string, msgHandler MsgHandler) (chan struct{}, error) {
	logrus.Infof("Binance, subscribe received on topic %s", topic)
	isExist := e.isSubscriptionExist(topic)
	if isExist == true {
		logrus.Tracef("Binance, subscription on topic %s already existed", topic)
		return nil, nil
	}
	elements := strings.Split(topic, "_")
	if len(elements) != 3 {
		return nil, errors.New("topic format incorrect (sample: binanace_kline_ethusdt:1m)")
	}
	ex := elements[0]
	if ex != "binance" {
		return nil, errors.New("incorrect exchange")
	}
	op := elements[1]
	switch strings.ToLower(op) {
	case "kline":
		stop, err := subscribeKline(topic, elements[2], msgHandler)
		if err != nil {
			return stop, err
		}
		e.subscriptions[topic] = stop
		return stop, nil
	case "depth":
		stop := subscribeDepth(topic, elements[2], msgHandler)
		e.subscriptions[topic] = stop
		return stop, nil
	default:
		return nil, errors.New("unknown topic")
	}
	return nil, nil
}

func (e *Exchange) isSubscriptionExist(topic string) bool {
	if _, ok := e.subscriptions[topic]; ok {
		logrus.Tracef("Subscription on %s already exists", topic)
		return true
	}
	return false
}

func subscribeKline(topic string, symbolAndInterval string, msgHandler MsgHandler) (stop chan struct{}, err error) {
	klineElements := strings.Split(symbolAndInterval, ":")
	if len(klineElements) != 2 {
		err = errors.New("kline format incorrect")
	}
	symbol := strings.ToUpper(klineElements[0])
	interval := strings.ToLower(klineElements[1])
	wsKlineHandler := func(event *bin.WsKlineEvent) {
		j, _ := json.Marshal(event)
		str := string(j)
		logrus.Infof("subscribeKline: %s", str)
		msgHandler(topic, str)
	}
	logrus.Infof("ListenKline on %s %s", symbol, interval)
	stop = binance.ListenKline(symbol, interval, wsKlineHandler)
	return stop, nil
}

func subscribeDepth(topic string, symbol string, msgHandler MsgHandler) (stop chan struct{}) {
	wsDepthHandler := func(event *bin.WsDepthEvent) {
		j, _ := json.Marshal(event)
		str := string(j)
		msgHandler(topic, str)
	}
	logrus.Infof("ListenDepth on %s", symbol)
	stop = binance.ListenDepth(symbol, wsDepthHandler)
	return stop
}
