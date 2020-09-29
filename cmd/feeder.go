//
//
//
package main

import (
	feeder "coinbase_feed"
	"coinbase_feed/models"
	"coinbase_feed/storage"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

var err error

type listSymbols []string

func (i *listSymbols) String() string {
	return fmt.Sprintf("symbols: %v", *i)
}

func (i *listSymbols) Set(value string) error {
	*i = append(*i, value)
	return nil
}

const (
	MYSQL_URL = "MYSQL_URL"
)

func main() {
	var symbols listSymbols
	flag.Var(&symbols, "symbol", "list of symbols - exp: -symbol ETH-BTC  -symbol BTC-USD")
	flag.Parse()

	if len(symbols) == 0 {
		panic("empty symbols list\n Use args: -symbol ETH-BTC  -symbol BTC-USD")
	}

	mysqlConnectionURL := os.Getenv(MYSQL_URL)
	if mysqlConnectionURL == "" {
		panic("empty mysql connection url. Use env MYSQL_URL")
	}

	// todo: parse symbols from flag.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())
	g, subCtx := errgroup.WithContext(ctx)

	// create storage
	client, err := storage.NewMySQLClient(mysqlConnectionURL)
	if err != nil {
		log.Fatal(err)
	}
	err = client.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}

	// make channels
	channels := models.NewTickerChannel(symbols)
	channels.Iter(func(symbol string, ch chan feeder.Ticker) {
		listeners(subCtx, g, symbol, client, ch)
	})
	tickersProducer(subCtx, g, channels)

	g.Go(func() error {
		sig := <-quit
		log.Println("Shutting down server... Reason:", sig)
		cancel()
		channels.CloseChannels()
		return client.Close()
	})

	err = g.Wait()
	if err != nil {
		log.Printf("Error wait %v \n", err)
	}
	log.Println("Ended")
}

func errHandler(err error) {
	log.Println("Errors:", err)
}

func tickersProducer(ctx context.Context, g *errgroup.Group, tickerChannels *models.TickerChannels) {
	g.Go(func() error {
	RECONNECT:
		opt := feeder.NewCoinbaseWsOptions()
		opt.Debug = true
		doneC, err := feeder.StartFeedTicker(ctx, opt, tickerChannels.GetSymbols(), func(t feeder.Ticker) {
			if ch, ok := tickerChannels.GetChannel(t.ProductID); ok {
				ch <- t
			}
		}, errHandler)
		if err != nil {
			return err
		}

		for {
			select {
			case <-doneC:
				if ctx.Err() == nil {
					log.Println("drop connection?")
					goto RECONNECT
				}
				return nil
			case <-ctx.Done():
				return nil
			}
		}
	})

}

func listeners(ctx context.Context, g *errgroup.Group, symbol string, client storage.Storager, input <-chan feeder.Ticker) {
	g.Go(func() error {
		log.Println("Start listener on ", symbol)

		for {
			select {
			case <-ctx.Done():
				// stop
				return nil

			case ticker := <-input:
				newTicker, err := models.ConvertTickerFromStreamTicker(ticker)
				if err != nil {
					errHandler(err)
					continue
				}
				fmt.Println(newTicker)
				err = client.Write(ctx, newTicker)
				if err != nil {
					errHandler(err)
					continue
				}

				// res, err := client.ReadLast(ctx, symbol)
				// fmt.Println(res, err)
			}
		}
	})
}
