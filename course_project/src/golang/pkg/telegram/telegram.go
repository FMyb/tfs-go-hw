package telegram

import (
	"context"
	"fmt"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type TelegramClient struct {
	token string
	repo  repositories.ClientRepository
	bot   *tgbotapi.BotAPI
}

func NewTelegramClient(token string, repo repositories.ClientRepository) *TelegramClient {
	return &TelegramClient{token: token, repo: repo}
}

func toMessage(order domain.ResponseOrder) string {
	return fmt.Sprintf("Send order: result: %s, orderID: %s, side: %s, price: %f, size: %f, symbol: %s, status: %s, server time: %s, type: %s",
		order.Result(),
		order.OrderId(),
		order.Side(),
		order.Price(),
		order.Quantity(),
		order.Symbol(),
		order.Status(),
		order.ServerTime(),
		order.Type(),
	)
}

func (tc *TelegramClient) SendOrder(order domain.ResponseOrder) error {
	log.Debug("start send order to client")
	chatIDs, err := tc.repo.Users(context.Background())
	if err != nil {
		return err
	}
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, toMessage(order))
		log.Debugf("send message %s to %d", msg.Text, chatID)
		_, err = tc.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tc *TelegramClient) Start() error {
	log.Info("Start tg bot")
	bot, err := tgbotapi.NewBotAPI(tc.token)
	if err != nil {
		return err
	}
	tc.bot = bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := tc.bot.GetUpdatesChan(updateConfig)
	go func() {
		log.Info("Start listen tg bot updates")
		for update := range updates {
			if update.Message == nil {
				continue
			}
			chatId := update.Message.Chat.ID
			log.Infof("New request from tg bot (chatID: %d, message: %s)", chatId, update.Message.Text)
			switch update.Message.Text {
			case "/start":
				err = tc.repo.SaveUser(context.Background(), chatId)
				if err != nil {
					continue
				}
			case "/stop":
				err = tc.repo.DeleteUser(context.Background(), chatId)
				if err != nil {
					continue
				}
			}
		}
	}()
	return nil
}

func (tc *TelegramClient) Stop() {
	tc.bot.StopReceivingUpdates()
	log.Info("Stop tg bot...")
}
