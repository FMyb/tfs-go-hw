package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestToSuccessResponseTicker(t *testing.T) {
	type args struct {
		ticker *ResponseTicker
	}
	tests := []struct {
		name string
		args args
		want SuccessResponseTicker
	}{
		{
			name: "simple to success response ticker",
			args: args{
				ticker: &ResponseTicker{
					Time:      42424242,
					Feed:      "ticker",
					ProductID: "PI_XBTUSD",
					Suspended: false,
					MarkPrice: 123123,
					Event:     "",
					Message:   "",
				},
			},
			want: SuccessResponseTicker{
				Time:      42424242,
				Feed:      "ticker",
				ProductID: "PI_XBTUSD",
				Suspended: false,
				MarkPrice: 123123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSuccessResponseTicker(tt.args.ticker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSuccessResponseTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorResponseTicker_Status(t *testing.T) {
	type fields struct {
		Event   string
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   Status
	}{
		{
			name: "error ticker status",
			fields: fields{
				Event:   "error",
				Message: "some_error",
			},
			want: ERROR,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrorResponseTicker{
				Event:   tt.fields.Event,
				Message: tt.fields.Message,
			}
			if got := e.Status(); got != tt.want {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_OrderId(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get order id",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId:     "123123123123asbas",
					Status:      "placed",
					OrderEvents: nil,
				},
			},
			want: "123123123123asbas",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.OrderId(); got != tt.want {
				t.Errorf("OrderId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Price(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "get price",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "",
								Side:   "",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: 1234,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Price(); got != tt.want {
				t.Errorf("Price() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Quantity(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "get quantity",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "",
								Side:   "",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: 1,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Quantity(); got != tt.want {
				t.Errorf("Quantity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Result(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   Status
	}{
		{
			name: "get result",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId:     "",
					Status:      "placed",
					OrderEvents: nil,
				},
			},
			want: SUCCESSES,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Result(); got != tt.want {
				t.Errorf("Result() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_ServerTime(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get server time",
			fields: fields{
				KserverTime: time.UnixMilli(123123123),
			},
			want: time.UnixMilli(123123123),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.ServerTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Side(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get side",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "",
								Side:   "buy",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: "buy",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Side(); got != tt.want {
				t.Errorf("Side() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Status(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get status",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "",
								Side:   "",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: "placed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Status(); got != tt.want {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Symbol(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get quantity",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "PI_XBTUSD",
								Side:   "",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: "PI_XBTUSD",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Symbol(); got != tt.want {
				t.Errorf("Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKrakenResponseOrder_Type(t *testing.T) {
	type fields struct {
		Kresult     Status
		Kerror      string
		KserverTime time.Time
		sendStatus  SendStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   Type
	}{
		{
			name: "get quantity",
			fields: fields{
				Kresult:     SUCCESSES,
				Kerror:      "",
				KserverTime: time.Time{},
				sendStatus: SendStatus{
					OrderId: "",
					Status:  "placed",
					OrderEvents: []OrderEvent{{
						ExecutionOrderEvent: ExecutionOrderEvent{
							ExecPrice:  1234,
							ExecAmount: 1,
							OrderPriorExecution: struct {
								Symbol string `json:"symbol"`
								Side   string `json:"side"`
							}{
								Symbol: "",
								Side:   "",
							},
						},
						OType: EXECUTION,
					}},
				},
			},
			want: EXECUTION,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KrakenResponseOrder{
				Kresult:     tt.fields.Kresult,
				Kerror:      tt.fields.Kerror,
				KserverTime: tt.fields.KserverTime,
				SendStatus:  tt.fields.sendStatus,
			}
			if got := k.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuccessResponseTicker_Status(t *testing.T) {
	type fields struct {
		Time      uint64
		Feed      string
		ProductID string
		Suspended bool
		MarkPrice float32
	}
	tests := []struct {
		name   string
		fields fields
		want   Status
	}{
		// TODO: Add test cases.
		{
			name: "success status",
			fields: fields{
				Time:      0,
				Feed:      "",
				ProductID: "",
				Suspended: false,
				MarkPrice: 0,
			},
			want: SUCCESSES,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SuccessResponseTicker{
				Time:      tt.fields.Time,
				Feed:      tt.fields.Feed,
				ProductID: tt.fields.ProductID,
				Suspended: tt.fields.Suspended,
				MarkPrice: tt.fields.MarkPrice,
			}
			if got := s.Status(); got != tt.want {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderEvent_Price(t *testing.T) {
	type fields struct {
		ExecutionOrderEvent ExecutionOrderEvent
		OType               Type
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "event price",
			fields: fields{
				ExecutionOrderEvent: ExecutionOrderEvent{
					ExecPrice:  123123,
					ExecAmount: 0,
					OrderPriorExecution: struct {
						Symbol string `json:"symbol"`
						Side   string `json:"side"`
					}{
						Symbol: "",
						Side:   "",
					},
				},
				OType: EXECUTION,
			},
			want: 123123,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := OrderEvent{
				ExecutionOrderEvent: tt.fields.ExecutionOrderEvent,
				OType:               tt.fields.OType,
			}
			if got := e.Price(); got != tt.want {
				t.Errorf("Price() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderEvent_Quantity(t *testing.T) {
	type fields struct {
		ExecutionOrderEvent ExecutionOrderEvent
		OType               Type
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "event quantity",
			fields: fields{
				ExecutionOrderEvent: ExecutionOrderEvent{
					ExecPrice:  0,
					ExecAmount: 1,
					OrderPriorExecution: struct {
						Symbol string `json:"symbol"`
						Side   string `json:"side"`
					}{
						Symbol: "",
						Side:   "",
					},
				},
				OType: EXECUTION,
			},
			want: 1,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := OrderEvent{
				ExecutionOrderEvent: tt.fields.ExecutionOrderEvent,
				OType:               tt.fields.OType,
			}
			if got := e.Quantity(); got != tt.want {
				t.Errorf("Quantity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderEvent_Side(t *testing.T) {
	type fields struct {
		ExecutionOrderEvent ExecutionOrderEvent
		OType               Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "event side",
			fields: fields{
				ExecutionOrderEvent: ExecutionOrderEvent{
					ExecPrice:  0,
					ExecAmount: 0,
					OrderPriorExecution: struct {
						Symbol string `json:"symbol"`
						Side   string `json:"side"`
					}{
						Symbol: "PI_XBTUSD",
						Side:   "sell",
					},
				},
				OType: EXECUTION,
			},
			want: "sell",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := OrderEvent{
				ExecutionOrderEvent: tt.fields.ExecutionOrderEvent,
				OType:               tt.fields.OType,
			}
			if got := e.Side(); got != tt.want {
				t.Errorf("Side() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderEvent_Symbol(t *testing.T) {
	type fields struct {
		ExecutionOrderEvent ExecutionOrderEvent
		OType               Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "event symbol",
			fields: fields{
				ExecutionOrderEvent: ExecutionOrderEvent{
					ExecPrice:  0,
					ExecAmount: 0,
					OrderPriorExecution: struct {
						Symbol string `json:"symbol"`
						Side   string `json:"side"`
					}{
						Symbol: "PI_XBTUSD",
						Side:   "",
					},
				},
				OType: EXECUTION,
			},
			want: "PI_XBTUSD",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := OrderEvent{
				ExecutionOrderEvent: tt.fields.ExecutionOrderEvent,
				OType:               tt.fields.OType,
			}
			if got := e.Symbol(); got != tt.want {
				t.Errorf("Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
