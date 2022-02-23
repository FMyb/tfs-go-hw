package services

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/orders"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/pkg/telegram"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/tickers"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestRunBot(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	ctx, cancel := context.WithCancel(context.Background())
	type args struct {
		bot Bot
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple run bot",
			args: args{
				bot: NewMockBot(
					tickers.NewMockwsTicker(),
					repositories.NewMockOrderRepository(),
					ctx,
					telegram.NewMockClient(),
					orders.NewMockOrder(),
					cancel,
				),
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunBot(tt.args.bot); (err != nil) != tt.wantErr {
				t.Errorf("RunBot() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(10 * time.Second)
			StopBot(tt.args.bot, tt.args.bot.Cancel())
		})
	}
}
