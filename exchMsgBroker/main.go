package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"hamednikzad/exchange-message-broker/pkg/bridge"
	"hamednikzad/exchange-message-broker/pkg/broker"
	"hamednikzad/exchange-message-broker/pkg/exchange"
	"sync"
)

var addr = flag.String("addr", ":8080", "http service address")
var wg = &sync.WaitGroup{}
var appSettings AppSettings

func init() {
	appSettings = getAppSettings()

	lvl := appSettings.Logging.Level

	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		fmt.Printf("Could not parse log level %s. so loglevel set to Debug\n", lvl)
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)

	logrus.Trace("Will only be visible if the loglevel permits it")
	logrus.Debug("Will only be visible if the loglevel permits it")
	logrus.Info("Will only be visible if the loglevel permits it")
	logrus.Error("Will only be visible if the loglevel permits it")
	logrus.Warn("Will only be visible if the loglevel permits it")

	addr = &appSettings.Server.Listen
}

func main() {

	flag.Parse()

	wg.Add(1)

	hub := broker.Run(*addr)
	ex := exchange.NewExchange()
	var brg = bridge.NewBridge(
		hub,
		ex,
	)
	go brg.Run()

	wg.Wait()
}

func handySub() {
	//stopC := make(chan struct{})
	/*wg.Add(1)

	stopKline, _ := ex.Subscribe("binance_kline_btcusdt:1m", nil)
	stopDepth, _ := ex.Subscribe("binance_depth_btcusdt", nil)

		//depthFinished := make(chan struct{})
		//go exchange.Subscribe("binance_depth_btcusdt", nil, depthFinished)
	go func() {
		//defer wg.Done()
		time.Sleep(5 * time.Second)
		logrus.Trace("Sleep over")
		ex.UnSubscribe("binance_kline_btcusdt:1m")
		ex.UnSubscribe("binance_depth_btcusdt")

		wg.Done()
		return
		stopKline <- struct{}{}
		stopDepth <- struct{}{}
	}()
	*/
	//<-depthFinished
}
