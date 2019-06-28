package main

import (
	"flag"

	"github.com/Kucoin/kucoin-go-sdk"
	order_book "github.com/mrlooper/kucoin-go-level2-demo"
	"github.com/mrlooper/kucoin-go-level2-demo/web"
)

func main() {
	symbol, port := getArgs()

	apiService := kucoin.NewApiServiceFromEnv()
	l2OrderBook := order_book.NewLevel2OrderBook(apiService, symbol)
	go l2OrderBook.ReloadOrderBook()

	r := web.NewRouter(port, l2OrderBook)
	go r.Handle()

	websocket(apiService, l2OrderBook)
}

func getArgs() (string, string) {
	symbol := flag.String("s", "BTC-USDT", "symbol")
	port := flag.String("p", "9090", "port")
	flag.Parse()

	return *symbol, *port
}

func websocket(apiService *kucoin.ApiService, l2OrderBook *order_book.Level2OrderBook) {
	rsp, err := apiService.WebSocketPublicToken()
	if err != nil {
		panic(err)
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		panic(err)
	}

	c := apiService.NewWebSocketClient(tk)

	mc, ec, err := c.Connect()
	if err != nil {
		panic(err)
	}

	ch := kucoin.NewSubscribeMessage("/market/level2:"+l2OrderBook.Symbol(), false)
	if err := c.Subscribe(ch); err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-ec:
			c.Stop() // Stop subscribing the WebSocket feed
			panic(err)

		case msg := <-mc:
			l2OrderBook.Messages <- msg.RawData
		}
	}
}
