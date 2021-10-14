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
		candle := domain.Candle{}
		for price := range in {
			time1m, err := domain.PeriodTS(domain.CandlePeriod1m, price.TS)
			if errors.Is(err, domain.ErrUnknownPeriod) {
				log.Errorf("Can't convert period %s", price.TS)
				continue
			}
			ok := false
			candle, ok = candles[price.Ticker]
			if ok {
				candle.Low = math.Min(candle.Low, price.Value)
				candle.High = math.Max(candle.High, price.Value)
				candle.Close = price.Value
				if candle.TS == time1m {
					candles[price.Ticker] = candle
				} else {
					out <- candle
					candles[price.Ticker] = createCandle1mFromPrice(price, time1m)
					candle = candles[price.Ticker]
				}
			} else {
				candles[price.Ticker] = createCandle1mFromPrice(price, time1m)
				candle = candles[price.Ticker]
			}
		}
		out <- candle
		log.Infof("Close 1m")
	}()
	return out
}

func toOtherTime(in <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	candles := make(map[string]domain.Candle)
	go func() {
		defer close(out)
		candle := domain.Candle{}
		for candleM := range in {
			ts, err := domain.PeriodTS(period, candleM.TS)
			if errors.Is(err, domain.ErrUnknownPeriod) {
				log.Errorf("Can't convert period %s", candleM.TS)
				continue
			}
			ok := false
			candle, ok = candles[candleM.Ticker]
			if ok {
				candle.Low = math.Min(candle.Low, candleM.Low)
				candle.High = math.Max(candle.High, candleM.High)
				candle.Close = candleM.Close
				if candle.TS == ts {
					candles[candleM.Ticker] = candle
				} else {
					out <- candle
					candles[candleM.Ticker] = createCandleFromCandle(candleM, period, ts)
					candle = candles[candleM.Ticker]
				}
			} else {
				candles[candleM.Ticker] = createCandleFromCandle(candleM, period, ts)
				candle = candles[candleM.Ticker]
			}
		}
		out <- candle
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

func writeCandlesToFile(in <-chan domain.Candle, fileName string, wg *sync.WaitGroup) <-chan domain.Candle {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Can't create candles_1m.csv: %s", err)
		os.Exit(0)
	}
	out := make(chan domain.Candle)
	go func() {
		defer wg.Done()
		defer close(out)
		w := csv.NewWriter(f)
		defer w.Flush()
		i := 0
		for candle := range in {
			log.Infof("prices (%s) %d: %+v", candle.Period, i, candle)
			err := w.Write(candleToString(candle))
			if err != nil {
				log.Errorf("Can't write candle: %s", err)
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
	wg := sync.WaitGroup{}
	wg.Add(4)
	prices := pg.Prices(ctx)
	candles1m := writeCandlesToFile(to1m(prices), "./lection03/homework/candles_1m.csv", &wg)
	candles2m := writeCandlesToFile(to2m(candles1m), "./lection03/homework/candles_2m.csv", &wg)
	candles10m := writeCandlesToFile(to10m(candles2m), "./lection03/homework/candles_10m.csv", &wg)
	go func() {
		defer wg.Done()
		for i := range candles10m {
			logger.Infof("Candle processing is completed: %+v", i)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	logger.Infof("Start shoutdown...")
	cancel()
	wg.Wait()
}
