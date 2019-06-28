package order_book

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/mrlooper/kucoin-go-level2-demo/helper"
	"github.com/Kucoin/kucoin-go-sdk"
)

type Level2OrderBook struct {
	apiService *kucoin.ApiService
	symbol     string
	lock       *sync.RWMutex
	Messages   chan json.RawMessage

	partOrderBook *PartOrderBookModel
}

func NewLevel2OrderBook(apiService *kucoin.ApiService, symbol string) *Level2OrderBook {
	l2book := &Level2OrderBook{
		apiService: apiService,
		symbol:     symbol,
		lock:       &sync.RWMutex{},
		Messages:   make(chan json.RawMessage, 300),
	}

	return l2book
}

func (l2book *Level2OrderBook) Symbol() string {
	return l2book.symbol
}

func (l2book *Level2OrderBook) resetOrderBook() {
	l2book.lock.Lock()
	l2book.partOrderBook = &PartOrderBookModel{}
	l2book.lock.Unlock()
}

func (l2book *Level2OrderBook) ReloadOrderBook() {
	defer func() {
		if r := recover(); r != nil {
			helper.Error("ReloadOrderBook panic: %v", r)
			l2book.ReloadOrderBook()
		}
	}()

	helper.Debug("Symbol: %s start ReloadOrderBook", l2book.symbol)
	l2book.resetOrderBook()

	l2book.playback()

	for msg := range l2book.Messages {
		helper.Debug("Proccesing message: " + kucoin.ToJsonString(msg))
		l2Data, err := NewLevel2StreamDataL2UpdateModel(msg)
		if err != nil {
			panic(err)
		}
		l2book.updateFromStream(l2Data)
	}
}

func (l2book *Level2OrderBook) playback() {
	helper.Warn("Prepare playback...")

	partOrderBook, err := l2book.getAggregatedPartOrderBook()
	helper.Info("Loaded OrderBook from API. %v asks and %v bids", len(partOrderBook.Asks), len(partOrderBook.Bids))
	helper.Debug("OrderBook sequence: " + partOrderBook.Sequence)
	if err != nil {
		panic("Error loading order book")
	}

	l2book.lock.Lock()
	l2book.partOrderBook = partOrderBook
	l2book.lock.Unlock()

	helper.Debug("Playback finished")
	helper.Debug(fmt.Sprintf("Asks: %v", l2book.partOrderBook.Asks))
	helper.Debug(fmt.Sprintf("Bids: %v", l2book.partOrderBook.Bids))

}

func (l2book *Level2OrderBook) updateFromStream(msg *Level2StreamDataL2UpdateModel) {
	l2book.lock.Lock()
	defer l2book.lock.Unlock()

	l2book.updateOrderBook(msg)
	l2book.updateSequence(msg)

}

func (l2book *Level2OrderBook) updateSequence(msg *Level2StreamDataL2UpdateModel) (bool) {
	l2book.partOrderBook.Sequence = strconv.FormatUint(msg.SequenceEnd, 10)

	return true
}

func (l2book *Level2OrderBook) updateOrderBook(msg *Level2StreamDataL2UpdateModel) {
	partOrderBookSequence, err := helper.Uint64FromString(l2book.partOrderBook.Sequence)
	if err != nil {
		panic("format failed: " + l2book.partOrderBook.Sequence)
	}

	helper.Debug("Updating OrderBook. Sequence: " + l2book.partOrderBook.Sequence)
	n := len(msg.Changes.Asks)
	for i := 0; i < n; i++ {
		ask := msg.Changes.Asks[i]
		seq, err := helper.Uint64FromString(ask[2])
		if err != nil {
			panic("format failed: " + l2book.partOrderBook.Sequence)
		}
		if seq > partOrderBookSequence {
			l2book.partOrderBook.Asks = updateAskFromLevel2OrderBook(l2book.partOrderBook.Asks, ask[0], ask[1])
		} else{
			helper.Debug(fmt.Sprintf("Discarding ask %v", ask))
		}
	}

	n = len(msg.Changes.Bids)
	for i := 0; i < n; i++ {
		bid := msg.Changes.Bids[i]
		seq, err := helper.Uint64FromString(bid[2])
		if err != nil {
			panic("format failed: " + l2book.partOrderBook.Sequence)
		}
		if seq > partOrderBookSequence {
			l2book.partOrderBook.Bids = updateBidFromLevel2OrderBook(l2book.partOrderBook.Bids, bid[0], bid[1])
		} else{
			helper.Debug(fmt.Sprintf("Discarding bid %v", bid))
		}		
	}

	helper.Debug(fmt.Sprintf("Asks: %v", l2book.partOrderBook.Asks))
	helper.Debug(fmt.Sprintf("Bids: %v", l2book.partOrderBook.Bids))

}

func (l2book *Level2OrderBook) getAggregatedPartOrderBook() (*PartOrderBookModel, error) {
	rsp, err := l2book.apiService.AggregatedPartOrderBook(l2book.symbol, 100)
	if err != nil {
		return nil, err
	}

	c := &PartOrderBookModel{}
	if err := rsp.ReadData(c); err != nil {
		return nil, err
	}

	if c.Sequence == "" {
		return nil, errors.New("empty key sequence")
	}

	return c, nil
}

func (l2book *Level2OrderBook) SnapshotBytes() ([]byte, error) {
	l2book.lock.RLock()
	data, err := json.Marshal(*l2book.partOrderBook)
	l2book.lock.RUnlock()
	if err != nil {
		return nil, err
	}

	return data, nil
}
