package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fardog/primebot"

	"github.com/dghubble/oauth1"
)

var (
	interval = flag.Duration(
		"interval",
		1*time.Hour,
		"Interval at which primes should be posted.",
	)
	consumerKey = flag.String(
		"consumer-key",
		os.Getenv("TWITTER_CONSUMER_KEY"),
		"Consumer Key to use when connecting. Falls back to value in TWITTER_CONSUMER_KEY env var",
	)
	consumerSecret = flag.String(
		"consumer-secret",
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		"Consumer Secret to use when connecting. Falls back to value in TWITTER_CONSUMER_SECRET env var",
	)
	accessToken = flag.String(
		"access-token",
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		"Access Token Key to use when connecting. Falls back to value in TWITTER_ACCESS_TOKEN env var",
	)
	accessSecret = flag.String(
		"access-secret",
		os.Getenv("TWITTER_ACCESS_SECRET"),
		"Access Token Secret to use when posting. Falls back to value in TWITTER_ACCESS_SECRET env var",
	)
)

func main() {
	flag.Usage = func() {
		_, exe := filepath.Split(os.Args[0])
		fmt.Fprint(os.Stderr, "A Twitter implementation for primebot.")
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options]\n\nOptions:\n\n", exe)
		flag.PrintDefaults()
	}
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	client, err := primebot.NewTwitterClient(httpClient)
	if err != nil {
		log.Fatalf("error instantiating twitter client: %v", err)
	}

	t := primebot.NewIntervalTicker(*interval)
	g := primebot.NewProbablyPrimeGenerator(0)

	bot := primebot.NewPrimeBot(client, t, g, client, &primebot.BotOpts{
		Logger: logger,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		logger.Fatal(bot.Start(ctx))
	}()

	// run until exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel()
}
