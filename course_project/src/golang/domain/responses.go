package domain

import "time"

type ResponseOrder interface {
	Result() Status
	OrderId() string
	Status() string
	Symbol() string
	Quantity() float32
	Price() float32
	Type() Type
	ServerTime() time.Time
	Side() string
}

type Type = string

const (
	PLACE     Type = "PLACE"
	CANCEL    Type = "CANCEL"
	REJECT    Type = "REJECT"
	EDIT      Type = "EDIT"
	EXECUTION Type = "EXECUTION"
)

type KrakenResponseOrder struct {
	Kresult     Status    `json:"result,omitempty"`
	Kerror      string    `json:"error,omitempty"`
	KserverTime time.Time `json:"serverTime"`
	SendStatus  `json:"sendStatus"`
}

func (k KrakenResponseOrder) Result() Status {
	return k.Kresult
}

func (k KrakenResponseOrder) OrderId() string {
	return k.SendStatus.OrderId
}

func (k KrakenResponseOrder) Status() string {
	return k.SendStatus.Status
}

func (k KrakenResponseOrder) Symbol() string {
	return k.OrderEvents[0].Symbol()
}

func (k KrakenResponseOrder) Side() string {
	return k.OrderEvents[0].Side()
}

func (k KrakenResponseOrder) Quantity() float32 {
	return k.OrderEvents[0].Quantity()
}

func (k KrakenResponseOrder) Price() float32 {
	return k.OrderEvents[0].Price()
}

func (k KrakenResponseOrder) Type() Type {
	return k.OrderEvents[0].OType
}

func (k KrakenResponseOrder) ServerTime() time.Time {
	return k.KserverTime
}

type SendStatus struct {
	OrderId     string       `json:"order_id,omitempty"`
	Status      string       `json:"status,omitempty"`
	OrderEvents []OrderEvent `json:"orderEvents,omitempty"`
}

type OrderEvent struct {
	ExecutionOrderEvent `json:""`
	OType               Type `json:"type,omitempty"`
}

func (e OrderEvent) Price() float32 {
	switch e.OType {
	case EXECUTION: // TODO Add some types
		return e.ExecPrice
	default:
		return 0
	}
}

func (e OrderEvent) Quantity() float32 {
	switch e.OType {
	case EXECUTION: // TODO Add some types
		return e.ExecAmount
	default:
		return 0
	}
}

func (e OrderEvent) Symbol() string {
	switch e.OType {
	case EXECUTION:
		return e.OrderPriorExecution.Symbol
	default:
		return ""
	}
}

func (e OrderEvent) Side() string {
	switch e.OType {
	case EXECUTION:
		return e.OrderPriorExecution.Side
	default:
		return ""
	}
}

type ExecutionOrderEvent struct {
	ExecPrice           float32 `json:"price,omitempty"`
	ExecAmount          float32 `json:"amount,omitempty"`
	OrderPriorExecution struct {
		Symbol string `json:"symbol"`
		Side   string `json:"side"`
	} `json:"orderPriorExecution"`
}

type ResponseStatus interface {
	Status() Status
}

type Status = string

const (
	SUCCESSES Status = "success"
	ERROR     Status = "error"
)

type ResponseTicker struct {
	Time      uint64  `json:"time"`
	Feed      string  `json:"feed"`
	ProductID string  `json:"product_id"`
	Suspended bool    `json:"suspended"`
	MarkPrice float32 `json:"markPrice"`
	Event     string  `json:"event"`
	Message   string  `json:"message"`
}

func ToSuccessResponseTicker(ticker *ResponseTicker) SuccessResponseTicker {
	return SuccessResponseTicker{
		Time:      ticker.Time,
		Feed:      ticker.Feed,
		ProductID: ticker.ProductID,
		Suspended: ticker.Suspended,
		MarkPrice: ticker.MarkPrice,
	}
}

type SuccessResponseTicker struct {
	Time      uint64  `json:"time"`
	Feed      string  `json:"feed"`
	ProductID string  `json:"product_id"`
	Suspended bool    `json:"suspended"`
	MarkPrice float32 `json:"markPrice"`
}

func (s SuccessResponseTicker) Status() Status {
	return SUCCESSES
}

type ErrorResponseTicker struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

func (e ErrorResponseTicker) Status() Status {
	return ERROR
}
