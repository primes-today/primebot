package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fardog/primebot"
)

var (
	interval = flag.Duration(
		"interval",
		5*time.Second,
		"Interval at which primes should be posted",
	)
)

func main() {
	f := primebot.NewRandFetcher()
	t := primebot.NewIntervalTicker(*interval)
	r := &primebot.PlainFormat{}
	g := primebot.NewProbablyPrimeGenerator(0)
	p := primebot.NewStdoutPoster()

	bot := primebot.NewPrimeBot(f, t, g, r, p, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		bot.Start(ctx)
	}()

	// serve until exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel()
}
