package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
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
	flag.Usage = func() {
		_, exe := filepath.Split(os.Args[0])
		fmt.Fprint(os.Stderr, "A stdout generator for testing primebot.")
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options]\n\nOptions:\n\n", exe)
		flag.PrintDefaults()
	}
	flag.Parse()

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
