package services

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	log "github.com/sirupsen/logrus"
)

type Tickers interface {
	NextTicker() domain.ResponseStatus
	Start([]string)
	Stop() chan interface{}
}

type Bot interface {
	Tickers() Tickers
	ProductId() string
	PublicApiKey() string
	Repository() repositories.OrderRepository
	Context() context.Context
	Client() Client
	PrivateApiKey() string
	Sell(ticker domain.SuccessResponseTicker) bool
	Buy(ticker domain.SuccessResponseTicker) bool
	SetLastPrice(price float32)
	Orders() Orders
	Cancel() context.CancelFunc
}

type Client interface {
	SendOrder(order domain.ResponseOrder) error
}

type Orders interface {
	SendOrder(
		productId string,
		size uint64,
		side string,
		publicApiKey string,
		privateApiKey string,
	) (domain.ResponseOrder, error)
}

func RunBot(bot Bot) error {
	log.Info("run bot....")
	botTickers := bot.Tickers()
	botTickers.Start([]string{bot.ProductId()})
	resp, err := bot.Orders().SendOrder(bot.ProductId(), 1, "buy", bot.PublicApiKey(), bot.PrivateApiKey())
	bot.SetLastPrice(resp.Price())
	if err != nil {
		return err
	}
	err = bot.Repository().SendOrder(bot.Context(), resp)
	if err != nil {
		return err
	}
	err = bot.Client().SendOrder(resp)
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-bot.Context().Done():
		default:

		}
	stopFor:
		for {
			log.Info("start sell")
		stopSell:
			for {
				select {
				case <-bot.Context().Done():
					break stopSell
				default:
					switch t := botTickers.NextTicker().(type) {
					case domain.SuccessResponseTicker:
						switch {
						case bot.Sell(t):
							log.Info("send order to sell")
							resp, err := bot.Orders().
								SendOrder(bot.ProductId(), 1, "sell", bot.PublicApiKey(), bot.PrivateApiKey())
							if err != nil {
								log.Errorf("error in send order: %s", err)
								break stopFor
							}
							log.Debug("order sent")
							err = bot.Repository().SendOrder(bot.Context(), resp)
							if err != nil {
								log.Errorf("error in save order to database: %s", err)
								break stopFor
							}
							log.Debug("order wrote to repository")
							err = bot.Client().SendOrder(resp)
							if err != nil {
								log.Errorf("error in send order to client: %s", err)
								break stopFor
							}
							log.Debug("sent message to client")
							break stopSell
						default:
							continue
						}
					case domain.ErrorResponseTicker:
						log.Errorf("error in response ticker: %s", t.Message)
						continue
					default:
						log.Errorf("unsupported response ticker!!")
						return
					}
				}
			}
			log.Info("start buy")
		stopBuy:
			for {
				select {
				case <-bot.Context().Done():
					break stopFor
				default:
					switch t := botTickers.NextTicker().(type) {
					case domain.SuccessResponseTicker:
						switch {
						case bot.Buy(t):
							log.Info("send order to buy")
							resp, err := bot.Orders().SendOrder(bot.ProductId(), 1, "buy", bot.PublicApiKey(), bot.PrivateApiKey())
							if err != nil {
								log.Errorf("error in send order: %s", err)
								break stopFor
							}
							log.Debug("order sent")
							err = bot.Repository().SendOrder(bot.Context(), resp)
							if err != nil {
								log.Errorf("error in save order to database: %s", err)
								break stopFor
							}
							log.Debug("order wrote to repository")
							err = bot.Client().SendOrder(resp)
							if err != nil {
								log.Errorf("error in send order to client: %s", err)
								break stopFor
							}
							log.Debug("sent message to client")
							break stopBuy
						default:
							continue
						}
					case domain.ErrorResponseTicker:
						log.Errorf("error in response ticker: %s", t.Message)
						continue
					default:
						log.Errorf("unsupported response ticker!!")
						return
					}
				}
			}
		}
		<-botTickers.Stop()
		log.Info("stop running bot")
	}()
	return nil
}

func StopBot(bot Bot, cancelFunc context.CancelFunc) {
	log.Info("stop bot")
	cancelFunc()
}
