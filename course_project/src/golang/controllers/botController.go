package controllers

import (
	"context"
	"encoding/json"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/services"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

type Bot struct {
	BotProductId     string                       `json:"product_id,omitempty"`
	BotStopVal       *float32                     `json:"stop_val,omitempty"`
	Configured       bool                         `json:"configured,omitempty"`
	Started          bool                         `json:"started,omitempty"`
	BotTickers       services.Tickers             `json:"tickers,omitempty"`
	BotProfitVal     *float32                     `json:"profit_val,omitempty"`
	BotPublicApiKey  string                       `json:"public_api_key,omitempty"`
	BotRepository    repositories.OrderRepository `json:"repository,omitempty"`
	BotContext       context.Context              `json:"context,omitempty"`
	BotClient        services.Client              `json:"client,omitempty"`
	BotPrivateApiKey string                       `json:"private_api_key,omitempty"`
	BotOrders        services.Orders              `json:"bot_orders,omitempty"`
	mut              sync.Mutex
	lastPrice        float32
	cancel           context.CancelFunc
}

func NewBot(
	tickers services.Tickers,
	repository repositories.OrderRepository,
	ctx context.Context,
	client services.Client,
	orders services.Orders,
	publicApiKey string,
	privateApiKey string,
	cancel context.CancelFunc,
) *Bot {
	return &Bot{
		BotTickers:       tickers,
		BotRepository:    repository,
		BotContext:       ctx,
		BotClient:        client,
		BotOrders:        orders,
		BotPublicApiKey:  publicApiKey,
		BotPrivateApiKey: privateApiKey,
		cancel:           cancel,
	}
}

func NewConfiguredBot(
	productId string,
	stopVal float32,
	tickers services.Tickers,
	profitVal float32,
	publicApiKey string,
	privateApiKey string,
	repository repositories.OrderRepository,
	context context.Context,
	client services.Client,
	orders services.Orders,
	cancel context.CancelFunc,
) *Bot {
	return &Bot{
		BotProductId:     productId,
		BotStopVal:       &stopVal,
		BotTickers:       tickers,
		BotProfitVal:     &profitVal,
		BotPublicApiKey:  publicApiKey,
		BotRepository:    repository,
		BotContext:       context,
		BotClient:        client,
		BotPrivateApiKey: privateApiKey,
		Configured:       true,
		BotOrders:        orders,
		cancel:           cancel,
	}
}

func (b Bot) PrivateApiKey() string {
	return b.BotPrivateApiKey
}

func (b Bot) PublicApiKey() string {
	return b.BotPublicApiKey
}

func (b Bot) Repository() repositories.OrderRepository {
	return b.BotRepository
}

func (b Bot) Context() context.Context {
	return b.BotContext
}

func (b Bot) Client() services.Client {
	return b.BotClient
}

func (b Bot) Tickers() services.Tickers {
	return b.BotTickers
}

func (b Bot) ProductId() string {
	return b.BotProductId
}

func (b *Bot) SetLastPrice(price float32) {
	b.lastPrice = price
}

func (b Bot) Orders() services.Orders {
	return b.BotOrders
}

func (b Bot) Cancel() context.CancelFunc {
	return b.cancel
}

func (b *Bot) Sell(ticker domain.SuccessResponseTicker) bool {
	if b.lastPrice-ticker.MarkPrice >= *b.BotStopVal {
		return true
	}
	if ticker.MarkPrice-b.lastPrice >= *b.BotProfitVal {
		return true
	}
	return false
}

func (b *Bot) Buy(ticker domain.SuccessResponseTicker) bool {
	b.lastPrice = ticker.MarkPrice
	time.Sleep(5 * time.Second)
	return true
}

func (b *Bot) Configure(w http.ResponseWriter, r *http.Request) { // TODO parallel requests, вынести в сервис
	b.mut.Lock()
	defer b.mut.Unlock()
	log.Info("bot start configure")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error in read body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	bot := Bot{}
	err = json.Unmarshal(d, &bot)
	if err != nil || bot.BotProductId == "" || bot.BotStopVal == nil || b.Started {
		log.Info("bad request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b.configure(&bot)
	res, err := json.Marshal(bot)
	if err != nil {
		log.Errorf("error in marshal response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		log.Errorf("error in write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("bot finish configure success")
}

func (b *Bot) Start(w http.ResponseWriter, r *http.Request) {
	b.mut.Lock()
	defer b.mut.Unlock()
	log.Info("Bot starting...")
	defer r.Body.Close()
	if !b.Configured {
		log.Errorf("bot is not configured")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if b.Started {
		return
	}
	b.Started = true
	res, err := json.Marshal(b)
	if err != nil {
		log.Errorf("error in marshal request")
		w.WriteHeader(http.StatusInternalServerError)
		b.Started = false
		return
	}
	_, err = w.Write(res)
	if err != nil {
		log.Errorf("error in write result: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		b.Started = false
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	b.BotContext = ctx
	b.cancel = cancel
	err = services.RunBot(b)
	if err != nil {
		log.Errorf("error in start bot: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		b.Started = false
		return
	}
	log.Infof("Bot launched!")
}

func (b *Bot) Stop(w http.ResponseWriter, r *http.Request) {
	b.mut.Lock()
	defer r.Body.Close()
	if !b.Started {
		return
	}
	b.Started = false
	res, err := json.Marshal(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b.Started = true
		return
	}
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b.Started = true
		return
	}
	services.StopBot(b, b.cancel)
	b.mut.Unlock()
}

func (b *Bot) configure(bot *Bot) {
	b.BotProductId = bot.BotProductId
	b.BotStopVal = bot.BotStopVal
	b.BotProfitVal = bot.BotProfitVal
	b.Configured = true
}

func (b *Bot) Copy(bot *Bot) {
	b.BotProductId = bot.BotProductId
	b.BotStopVal = bot.BotStopVal
	b.Configured = bot.Configured
	b.Started = bot.Started
	b.BotTickers = bot.BotTickers
	b.BotProfitVal = bot.BotProfitVal
	b.BotPublicApiKey = bot.BotPublicApiKey
	b.BotRepository = bot.BotRepository
	b.BotContext = bot.BotContext
	b.BotClient = bot.BotClient
	b.BotPrivateApiKey = bot.BotPrivateApiKey
	b.lastPrice = bot.lastPrice
}
