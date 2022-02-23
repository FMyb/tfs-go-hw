package main

import (
	"context"
	"fmt"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/controllers"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/orders"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/pkg/postgres"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/pkg/telegram"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/tickers"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/utils"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	LogLevel       = "log.level"
	PublicApiKey   = "apikey.public"
	PrivateApiKey  = "apikey.private"
	TelegramToken  = "telegram.token"
	PostgresqlHost = "postgresql.host"
	PostgresqlMode = "postgresql.sslmode"
)

func main() {
	resourceManager := utils.NewResourceManager([]string{"course_project/resources/application.properties"})
	loglevel := resourceManager.GetProp(LogLevel)
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
	interrupt := make(chan interface{})
	wsTickers := tickers.NewWebSocketTicker(interrupt)
	publicApiKey := resourceManager.GetProp(PublicApiKey)
	if publicApiKey == "" {
		log.Fatal("Public api key must be in properties file")
	}
	privateApiKey := resourceManager.GetProp(PrivateApiKey)
	if privateApiKey == "" {
		log.Fatal("Private api key must be in properties file")
	}
	token := resourceManager.GetProp(TelegramToken)
	if token == "" {
		log.Fatal("Telegram bot token must be in properties file")
	}
	postgreHost := resourceManager.GetProp(PostgresqlHost)
	if postgreHost == "" {
		log.Fatal("PostgreSQL host must be in properties file")
	}
	sslmode := resourceManager.GetProp(PostgresqlMode)
	if sslmode == "" {
		sslmode = "disable"
	}
	pool, err := postgres.NewPool(fmt.Sprintf("postgresql://%s?sslmode=%s", postgreHost, sslmode))
	if err != nil {
		return
	}
	orderRepo := repositories.NewOrderRepository(pool)
	ctx, cancel := context.WithCancel(context.Background())
	telegramRepo := repositories.NewTelegramRepository(pool)
	client := telegram.NewTelegramClient(token, telegramRepo)
	BotOrders := orders.NewDefaultKrakenOrders()
	bot := controllers.NewBot(
		wsTickers,
		orderRepo,
		ctx,
		client,
		BotOrders,
		publicApiKey,
		privateApiKey,
		cancel,
	)
	err = client.Start()
	if err != nil {
		log.Errorf("error in starting bot: %s", err)
	}

	root := chi.NewRouter()
	root.Use(middleware.Logger)

	root.Post("/api/v1/configure", bot.Configure)
	root.Post("/api/v1/start", bot.Start)
	root.Post("/api/v1/stop", bot.Stop)

	log.Fatal(http.ListenAndServe(":5000", root))
}
