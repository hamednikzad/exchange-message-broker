package binance

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
)

func ListenKline(symbol string, interval string, wsKlineHandler binance.WsKlineHandler) (stop chan struct{}) {
	/*wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event)
	}*/
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, stop, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	return stop
}

func ListenKline2(finished chan struct{}) {
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsKlineServe("LTCBTC", "1m", wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	d := <-doneC
	finished <- d
}

func ListenDepth(symbol string, wsDepthHandler binance.WsDepthHandler) (stop chan struct{}) {
	/*wsDepthHandler := func(event *binance.WsDepthEvent) {
		fmt.Println(event)
	}*/
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, stop, err := binance.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	return stop
}
