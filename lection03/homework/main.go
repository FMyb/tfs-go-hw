package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/FMyb/tfs-go-hw/lection03/homework/domain"

	"github.com/FMyb/tfs-go-hw/lection03/homework/generator"
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

func to1m(in <-chan domain.Price) <-chan domain.Candle {
	out := make(chan domain.Candle)
	candles := make(map[string]domain.Candle)
	go func() {
		defer close(out)
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
		}
		log.Infof("Close 1m")
	}()
	return out
}

func toOtherTime(in <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	candles := make(map[string]domain.Candle)
	go func() {
		defer close(out)
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
		}
		log.Infof("Close %s", period)
	}()
	return out
}

func to2m(in <-chan domain.Candle) <-chan domain.Candle {
	return toOtherTime(in, domain.CandlePeriod2m)
}

func to10m(in <-chan domain.Candle) <-chan domain.Candle {
	return toOtherTime(in, domain.CandlePeriod10m)
}

func writeCandlesToFile(in <-chan domain.Candle, f *os.File, wg *sync.WaitGroup) <-chan domain.Candle {
	out := make(chan domain.Candle)
	go func() {
		defer wg.Done()
		defer close(out)
		w := csv.NewWriter(f)
		i := 0
		for candle := range in {
			log.Infof("prices (1m) %d: %+v", i, candle)
			err := w.Write(candleToString(candle))
			if err != nil {
				log.Errorf("Can't write candle: %s", err)
			} else {
				w.Flush()
			}
			i++
			out <- candle
		}
	}()
	return out
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())
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
	wg := sync.WaitGroup{}
	wg.Add(5)
	prices := pg.Prices(ctx)
	candles1m := writeCandlesToFile(to1m(prices), f1m, &wg)
	candles2m := writeCandlesToFile(to2m(candles1m), f2m, &wg)
	candles10m := writeCandlesToFile(to10m(candles2m), f10m, &wg)
	go func() {
		defer wg.Done()
		for i := range candles10m {
			logger.Infof("Candle processing is completed: %+v", i)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		defer wg.Done()
		<-c
		logger.Infof("Start shoutdown...")
		cancel()
	}()
	wg.Wait()
}
