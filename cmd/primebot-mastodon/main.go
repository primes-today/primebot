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
)

var (
	interval = flag.Duration(
		"interval",
		1*time.Hour,
		"Interval at which primes should be posted.",
	)
	serviceTimeout = flag.Duration(
		"service-timeout",
		30*time.Second,
		"Timeouts used when communicating with the Mastodon server.",
	)
	server = flag.String(
		"server",
		os.Getenv("MASTODON_SERVER"),
		"Mastodon instance URL to connect to. Falls back to value in MASTODON_SERVER env var",
	)
	clientID = flag.String(
		"client-id",
		os.Getenv("MASTODON_CLIENT_ID"),
		"ClientID to use when connecting. Falls back to value in MASTODON_CLIENT_ID env var",
	)
	clientSecret = flag.String(
		"client-secret",
		os.Getenv("MASTODON_CLIENT_SECRET"),
		"ClientSecret to use when connecting. Falls back to value in MASTODON_CLIENT_SECRET env var",
	)
	accessToken = flag.String(
		"access-token",
		os.Getenv("MASTODON_ACCESS_TOKEN"),
		"AccessToken to use when connecting. Falls back to value in MASTODON_ACCESS_TOKEN env var",
	)
	accountID = flag.String(
		"account-id",
		os.Getenv("MASTODON_ACCOUNT_ID"),
		"AccountID to use when posting. Falls back to value in MASTODON_ACCOUNT_ID env var",
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

	logger := log.New(os.Stdout, "", 0)

	config := &primebot.MastodonConfig{
		Server:       *server,
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		AccessToken:  *accessToken,
		AccountID:    *accountID,
	}
	cctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := primebot.NewMastodonClient(cctx, config)
	cancel()
	if err != nil {
		log.Fatalf("error instantiating mastodon client: %v", err)
	}

	t := primebot.NewIntervalTicker(*interval)
	g := primebot.NewProbablyPrimeGenerator(0)

	bot := primebot.NewPrimeBot(client, t, g, client, &primebot.BotOpts{
		Logger:         logger,
		ServiceTimeout: *serviceTimeout,
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
