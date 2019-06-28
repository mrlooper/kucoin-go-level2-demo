# KuCoin Level2 Demo

Level 2 order book implementation based on [kucoin-go-level3-demo](https://github.com/Kucoin/kucoin-go-level3-demo)

### Dependencies

```
go get github.com/Kucoin/kucoin-go-sdk
go get github.com/shopspring/decimal
go get github.com/sirupsen/logrus
```


### Usage

```
go run cmd/main.go -s BTC-USDT -p 8080
```

or you can download the latest available release.

when you see ***Loaded OrderBook from API*** at the command line, open [http://localhost:8080/](http://localhost:8080/)

