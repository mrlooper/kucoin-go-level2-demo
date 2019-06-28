package order_book

type (
	//[2]string{price", "size"}
	Level2OrderBookItem [2]string

	Level2AsksOrderBook [][2]string
	Level2BidsOrderBook [][2]string
)

type PartOrderBookModel struct {
	Sequence string              `json:"sequence"`
	Asks     Level2AsksOrderBook `json:"asks"` //ask Sort price from low to high
	Bids     Level2BidsOrderBook `json:"bids"` //bid Sort price from high to low
}