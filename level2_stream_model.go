package order_book

import (
	"encoding/json"

)

//Level 2 websocket stream data
type Level2StreamDataModel struct {
	Type       string `json:"type"`
	Topic      string `json:"topic"`
	Subject    string `json:"subject"`
	rawMessage json.RawMessage
}

func NewLevel2StreamDataModel(msgData json.RawMessage) (*Level2StreamDataModel, error) {
	l2Data := &Level2StreamDataModel{}

	if err := json.Unmarshal(msgData, l2Data); err != nil {
		return nil, err
	}
	l2Data.rawMessage = msgData

	return l2Data, nil
}

func (l2Data *Level2StreamDataModel) GetRawMessage() json.RawMessage {
	return l2Data.rawMessage
}

type Level2StreamDataL2UpdateModel struct {
	SequenceStart    uint64 `json:"sequenceStart"`
	SequenceEnd      uint64 `json:"sequenceEnd"`
	Symbol           string `json:"symbol"`
	Changes          Level2StreamDataL2UpdateChangesModel `json:"changes"`
	rawMessage json.RawMessage
}

type Level2StreamDataL2UpdateChangesModel struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

func NewLevel2StreamDataL2UpdateModel(msgData json.RawMessage) (*Level2StreamDataL2UpdateModel, error) {
	l2Data := &Level2StreamDataL2UpdateModel{}

	if err := json.Unmarshal(msgData, l2Data); err != nil {
		return nil, err
	}
	l2Data.rawMessage = msgData

	return l2Data, nil
}

func (l2Data *Level2StreamDataL2UpdateModel) GetRawMessage() json.RawMessage {
	return l2Data.rawMessage
}

const (
	BuySide  = "buy"
	SellSide = "sell"

	LimitOrderType  = "limit"
	MarketOrderType = "market"

	Level2MessageDoneCanceled = "canceled"
	Level2MessageDoneFilled   = "filled"

	Level2MessageReceivedType = "received"
	Level2MessageOpenType     = "open"
	Level2MessageDoneType     = "done"
	Level2MessageMatchType    = "match"
	Level2MessageChangeType   = "change"
)


