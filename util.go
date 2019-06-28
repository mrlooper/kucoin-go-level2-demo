package order_book

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
	"github.com/mrlooper/kucoin-go-level2-demo/helper"
)

//[2]string{"price", "size"}
//ask Sort price from low to high
func insertLevel2AskOrderBook(data Level2AsksOrderBook, el Level2OrderBookItem) Level2AsksOrderBook {
	index := sort.Search(len(data), func(i int) bool {
		aF, err := decimal.NewFromString(data[i][0])
		if err != nil {
			panic("format failed: " + data[i][0])
		}
		bF, _ := decimal.NewFromString(el[0])
		if err != nil {
			panic("format failed: " + el[0])
		}

		return aF.GreaterThan(bF)
	})

	data = append(data, [2]string{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}

//[2]string{"price", "size"}
//bid Sort price from high to low
func insertLevel2BidOrderBook(data Level2BidsOrderBook, el Level2OrderBookItem) Level2BidsOrderBook {
	index := sort.Search(len(data), func(i int) bool {
		aF, err := decimal.NewFromString(data[i][0])
		if err != nil {
			panic("format failed: " + data[i][0])
		}
		bF, _ := decimal.NewFromString(el[0])
		if err != nil {
			panic("format failed: " + el[0])
		}

		return aF.LessThan(bF)
	})

	helper.Debug(fmt.Sprintf("Inserting [%v] at index %v", el, index))
	data = append(data, [2]string{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}

//[3]string{"orderId", "price", "size"}
func deleteOrderFromLevel2OrderBook(data [][2]string, price string) [][2]string {
	for index, item := range data {
		if price == item[0] {
			helper.Debug(fmt.Sprintf("Removing %v from index %v", price, index))
			return append(data[:index], data[index+1:]...)
		}
	}

	return data
}

//[2]string{"price", "size"}
func updateBidFromLevel2OrderBook(data Level2BidsOrderBook, price string, size string) Level2BidsOrderBook {


	for index, item := range data {
		if size == "0" {
			return deleteOrderFromLevel2OrderBook(data, price)
		}

		if price == item[0] {
			helper.Debug(fmt.Sprintf("Updating bid %v to %v", price, size))
			ret := append(data[:index], [2]string{item[0], size})
			return append(ret, data[index+1:]...)
		} 
	}

	return insertLevel2BidOrderBook(data, Level2OrderBookItem{price, size})
}

func updateAskFromLevel2OrderBook(data Level2AsksOrderBook, price string, size string) Level2AsksOrderBook {
	for index, item := range data {
		if size == "0" {
			return deleteOrderFromLevel2OrderBook(data, price)
		}
			
		if price == item[0] {
			helper.Debug(fmt.Sprintf("Updating ask %v to %v", price, size))
			ret := append(data[:index], [2]string{item[0], size})
			return append(ret, data[index+1:]...)
		} 
	}

	return insertLevel2AskOrderBook(data, Level2OrderBookItem{price, size})
}
