package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/FMyb/tfs-go-hw/lection03/homework/domain"
	"github.com/FMyb/tfs-go-hw/lection03/homework/generator"
	"math"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func createCandle1mFromPrice(price domain.Price, ts time.Time) domain.Candle {
	return domain.Candle{
		Ticker: price.Ticker,
		Period: domain.CandlePeriod1m,
		Open:   price.Value,
		Low:    price.Value,
		High:   price.Value,
		TS:     ts,
	}
}

func createCandleFromCandle(candle domain.Candle, period domain.CandlePeriod, ts time.Time) domain.Candle {
	return domain.Candle{
		Ticker: candle.Ticker,
		Period: period,
		Open:   candle.Open,
		High:   candle.High,
		Low:    candle.Low,
		TS:     ts,
	}
}

func candleToString(candle domain.Candle) []string {
	res := []string{
		candle.Ticker,
		candle.TS.Format(time.RFC3339),
		fmt.Sprintf("%f", candle.Open),
		fmt.Sprintf("%f", candle.High),
		fmt.Sprintf("%f", candle.Low),
		fmt.Sprintf("%f", candle.Close),
	}
	return res
}

func to1m(ctx context.Context, in <-chan domain.Price) <-chan domain.Candle {
	out := make(chan domain.Candle)
	candles := make(map[string]domain.Candle)
	go func() {
		for price := range in {
			time1m, err := domain.PeriodTS(domain.CandlePeriod1m, price.TS)
			if errors.Is(err, domain.ErrUnknownPeriod) {
				log.Errorf("Can't convert period %s", price.TS)
				continue
			}
			candle, ok := candles[price.Ticker]
			if ok {
				candle.Low = math.Min(candle.Low, price.Value)
				candle.High = math.Max(candle.High, price.Value)
				if candle.TS == time1m {
					candles[price.Ticker] = candle
				} else {
					candle.Close = price.Value
					out <- candle
					candles[price.Ticker] = createCandle1mFromPrice(price, time1m)
				}
			} else {
				candles[price.Ticker] = createCandle1mFromPrice(price, time1m)
			}
			select {
			case <-ctx.Done():
				break
			default:
				continue
			}
		}
	}()
	return out
}

func toOtherTime(ctx context.Context, in <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	candles := make(map[string]domain.Candle)
	go func() {
		for candleM := range in {
			ts, err := domain.PeriodTS(period, candleM.TS)
			if errors.Is(err, domain.ErrUnknownPeriod) {
				log.Errorf("Can't convert period %s", candleM.TS)
				continue
			}
			candle, ok := candles[candleM.Ticker]
			if ok {
				candle.Low = math.Min(candle.Low, candleM.Low)
				candle.High = math.Max(candle.High, candleM.High)
				if candle.TS == ts {
					candles[candleM.Ticker] = candle
				} else {
					candle.Close = candleM.Close
					out <- candle
					candles[candleM.Ticker] = createCandleFromCandle(candleM, period, ts)
				}
			} else {
				candles[candleM.Ticker] = createCandleFromCandle(candleM, period, ts)
			}
			select {
			case <-ctx.Done():
				break
			default:
				continue
			}
		}
	}()
	return out
}

func to2m(ctx context.Context, in <-chan domain.Candle) <-chan domain.Candle {
	return toOtherTime(ctx, in, domain.CandlePeriod2m)
}

func to10m(ctx context.Context, in <-chan domain.Candle) <-chan domain.Candle {
	return toOtherTime(ctx, in, domain.CandlePeriod10m)
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(ctx1)
	ctx10, cancel10 := context.WithCancel(ctx2)
	c := make(chan os.Signal, 1)
	go func() {
		<-c
		cancel10()
		cancel2()
		cancel1()
		os.Exit(0)
	}()
	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	f1m, err := os.Create("./lection03/homework/candles_1m.csv")
	if err != nil {
		fmt.Printf("Can't create candles_1m.csv: %s", err)
		return
	}
	f2m, err := os.Create("./lection03/homework/candles_2m.csv")
	if err != nil {
		fmt.Printf("Can't create candles_2m.csv: %s", err)
		return
	}
	f10m, err := os.Create("./lection03/homework/candles_10m.csv")
	if err != nil {
		fmt.Printf("Can't create candles_10m.csv: %s", err)
		return
	}
	prices := pg.Prices(ctx)
	candles1m := to1m(ctx1, prices)
	rCandles1m := make(chan domain.Candle)
	candles2m := to2m(ctx2, rCandles1m)
	rCandles2m := make(chan domain.Candle)
	candles10m := to10m(ctx10, rCandles2m)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		w := csv.NewWriter(f1m)
		i := 0
		for candle := range candles1m {
			logger.Infof("prices (1m) %d: %+v", i, candle)
			rCandles1m <- candle
			err := w.Write(candleToString(candle))
			if err != nil {
				logger.Errorf("Can't write candle: %s", err)
			} else {
				w.Flush()
			}
			i++
		}
	}()
	go func() {
		defer wg.Done()
		w := csv.NewWriter(f2m)
		i := 0
		for candle := range candles2m {
			logger.Infof("prices (2m) %d: %+v", i, candle)
			rCandles2m <- candle
			err := w.Write(candleToString(candle))
			if err != nil {
				logger.Errorf("Can't write candle: %s", err)
			} else {
				w.Flush()
			}
			i++
		}
	}()
	go func() {
		defer wg.Done()
		w := csv.NewWriter(f10m)
		i := 0
		for candle := range candles10m {
			logger.Infof("prices (10m) %d: %+v", i, candle)
			err := w.Write(candleToString(candle))
			if err != nil {
				logger.Errorf("Can't write candle: %s", err)
			} else {
				w.Flush()
			}
			i++
		}
	}()
	wg.Wait()
	cancel()
}
