package controllers

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/services"
	"reflect"
	"testing"
)

func TestBot_Buy(t *testing.T) {
	type fields struct {
		BotProductId     string
		BotStopVal       *float32
		Configured       bool
		Started          bool
		BotTickers       services.Tickers
		BotProfitVal     *float32
		BotPublicApiKey  string
		BotRepository    repositories.OrderRepository
		BotContext       context.Context
		BotClient        services.Client
		BotPrivateApiKey string
		lastPrice        float32
	}
	type args struct {
		ticker domain.SuccessResponseTicker
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "simple buy",
			fields: fields{
				lastPrice: 0,
			},
			args: args{
				ticker: domain.SuccessResponseTicker{
					MarkPrice: 10000,
				},
			},
			want: true,
		},
		{
			name: "",
			fields: fields{
				lastPrice: 1000,
			},
			args: args{
				ticker: domain.SuccessResponseTicker{
					MarkPrice: 100,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				BotProductId:     tt.fields.BotProductId,
				BotStopVal:       tt.fields.BotStopVal,
				Configured:       tt.fields.Configured,
				Started:          tt.fields.Started,
				BotTickers:       tt.fields.BotTickers,
				BotProfitVal:     tt.fields.BotProfitVal,
				BotPublicApiKey:  tt.fields.BotPublicApiKey,
				BotRepository:    tt.fields.BotRepository,
				BotContext:       tt.fields.BotContext,
				BotClient:        tt.fields.BotClient,
				BotPrivateApiKey: tt.fields.BotPrivateApiKey,
				lastPrice:        tt.fields.lastPrice,
			}
			if got := b.Buy(tt.args.ticker); got != tt.want {
				t.Errorf("Buy() = %v, want %v", got, tt.want)
			}
			if _ = b.Buy(tt.args.ticker); b.lastPrice != tt.args.ticker.MarkPrice {
				t.Errorf("latPrice = %v, want %v", b.lastPrice, tt.args.ticker.MarkPrice)
			}
		})
	}
}

func TestBot_Sell(t *testing.T) {
	type fields struct {
		BotProductId     string
		BotStopVal       *float32
		Configured       bool
		Started          bool
		BotTickers       services.Tickers
		BotProfitVal     *float32
		BotPublicApiKey  string
		BotRepository    repositories.OrderRepository
		BotContext       context.Context
		BotClient        services.Client
		BotPrivateApiKey string
		lastPrice        float32
	}
	type args struct {
		ticker domain.SuccessResponseTicker
	}
	var val float32 = 100
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "sell stop lose",
			fields: fields{
				BotStopVal:   &val,
				BotProfitVal: &val,
				lastPrice:    150,
			},
			args: args{
				ticker: domain.SuccessResponseTicker{
					MarkPrice: 0,
				},
			},
			want: true,
		},
		{
			name: "sell take profit",
			fields: fields{
				BotStopVal:   &val,
				BotProfitVal: &val,
				lastPrice:    150,
			},
			args: args{
				ticker: domain.SuccessResponseTicker{
					MarkPrice: 300,
				},
			},
			want: true,
		},
		{
			name: "test sell stop lose",
			fields: fields{
				BotStopVal:   &val,
				BotProfitVal: &val,
				lastPrice:    150,
			},
			args: args{
				ticker: domain.SuccessResponseTicker{
					MarkPrice: 160,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				BotProductId:     tt.fields.BotProductId,
				BotStopVal:       tt.fields.BotStopVal,
				Configured:       tt.fields.Configured,
				Started:          tt.fields.Started,
				BotTickers:       tt.fields.BotTickers,
				BotProfitVal:     tt.fields.BotProfitVal,
				BotPublicApiKey:  tt.fields.BotPublicApiKey,
				BotRepository:    tt.fields.BotRepository,
				BotContext:       tt.fields.BotContext,
				BotClient:        tt.fields.BotClient,
				BotPrivateApiKey: tt.fields.BotPrivateApiKey,
				lastPrice:        tt.fields.lastPrice,
			}
			if got := b.Sell(tt.args.ticker); got != tt.want {
				t.Errorf("Sell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBot(t *testing.T) {
	type args struct {
		tickers       services.Tickers
		repository    repositories.OrderRepository
		ctx           context.Context
		client        services.Client
		orders        services.Orders
		publicApiKey  string
		privateApiKey string
		cancel        context.CancelFunc
	}
	tests := []struct {
		name string
		args args
		want *Bot
	}{
		{
			name: "",
			args: args{
				tickers:       nil,
				repository:    nil,
				ctx:           nil,
				client:        nil,
				orders:        nil,
				publicApiKey:  "",
				privateApiKey: "",
				cancel:        nil,
			},
			want: &Bot{
				BotProductId:     "",
				BotStopVal:       nil,
				Configured:       false,
				Started:          false,
				BotTickers:       nil,
				BotProfitVal:     nil,
				BotPublicApiKey:  "",
				BotRepository:    nil,
				BotContext:       nil,
				BotClient:        nil,
				BotPrivateApiKey: "",
				lastPrice:        0,
				BotOrders:        nil,
				cancel:           nil,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBot(
				tt.args.tickers,
				tt.args.repository,
				tt.args.ctx,
				tt.args.client,
				tt.args.orders,
				tt.args.publicApiKey,
				tt.args.privateApiKey,
				tt.args.cancel,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBot() = %v, want %v", got, tt.want)
			}
		})
	}
}
