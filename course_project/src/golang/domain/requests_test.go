package domain

import (
	"reflect"
	"testing"
)

func TestNewRequestWS(t *testing.T) {
	type args struct {
		event      string
		feed       string
		productIDs []string
	}
	tests := []struct {
		name string
		args args
		want *RequestWS
	}{
		{
			name: "request subscribe",
			args: args{
				event:      "subscribe",
				feed:       "ticker",
				productIDs: []string{"PI_XBTUSD"},
			},
			want: &RequestWS{
				Event:      "subscribe",
				Feed:       "ticker",
				ProductIDs: []string{"PI_XBTUSD"},
			},
		},
		{
			name: "request unsubscribe",
			args: args{
				event:      "unsubscribe",
				feed:       "ticker",
				productIDs: []string{"PI_XBTUSD"},
			},
			want: &RequestWS{
				Event:      "unsubscribe",
				Feed:       "ticker",
				ProductIDs: []string{"PI_XBTUSD"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequestWS(tt.args.event, tt.args.feed, tt.args.productIDs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestWS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeTickerRequestWS(t *testing.T) {
	type args struct {
		productIDs []string
	}
	tests := []struct {
		name string
		args args
		want *RequestWS
	}{
		{
			name: "subscribe ticker constructor",
			args: args{
				productIDs: []string{"PI_XBTUSD"},
			},
			want: &RequestWS{
				Event:      "subscribe",
				Feed:       "ticker",
				ProductIDs: []string{"PI_XBTUSD"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubscribeTickerRequestWS(tt.args.productIDs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubscribeTickerRequestWS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsubscribeTickerRequestWS(t *testing.T) {
	type args struct {
		productIDs []string
	}
	tests := []struct {
		name string
		args args
		want *RequestWS
	}{
		{
			name: "unsubscribe ticker constructor",
			args: args{
				productIDs: []string{"PI_XBTUSD"},
			},
			want: &RequestWS{
				Event:      "unsubscribe",
				Feed:       "ticker",
				ProductIDs: []string{"PI_XBTUSD"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnsubscribeTickerRequestWS(tt.args.productIDs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnsubscribeTickerRequestWS() = %v, want %v", got, tt.want)
			}
		})
	}
}
