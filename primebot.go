package primebot

import (
	"context"
	"log"
	"time"
)

type noopWriter struct{}

func (w *noopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

type Bot interface {
	Start(context.Context) error
}

type BotOpts struct {
	Logger *log.Logger
}

func NewPrimeBot(f Fetcher, t Ticker, g Generator, p Poster, opts *BotOpts) *PrimeBot {
	if opts == nil {
		opts = &BotOpts{}
	}

	if opts.Logger == nil {
		opts.Logger = log.New(&noopWriter{}, "", 0)
	}

	return &PrimeBot{
		ftc:  f,
		tck:  t,
		gen:  g,
		pst:  p,
		opts: opts,
		log:  opts.Logger,
	}
}

type PrimeBot struct {
	ftc  Fetcher
	tck  Ticker
	gen  Generator
	pst  Poster
	opts *BotOpts
	log  *log.Logger
}

func (p *PrimeBot) Start(ctx context.Context) error {
	p.log.Print("fetching initial list of statuses")

	fetchctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := p.ftc.Fetch(fetchctx)
	if err != nil {
		return err
	}

	st := cur.ts.Sub(time.Now())

	next := time.Now().Add(st)
	p.log.Printf("retrieved status %v; next post at %v", cur, next)

	p.gen.SetStart(cur.num + 1)
	t := p.tck.Start(ctx, st)
	pc := make(chan uint64, 1)
	er := make(chan error)
	go func() {
		for {
			n, err := p.gen.Generate(ctx)
			if err != nil {
				er <- err
				break
			}
			pc <- n
		}
	}()

	for {
		select {
		case <-t:
			status := <-pc
			err := p.pst.Post(ctx, status)
			if err != nil {
				return err
			}
			p.log.Printf("posted status %v", status)
		case err := <-er:
			p.log.Printf("error received from generator, %v. stopping", err)
			return err
		case <-ctx.Done():
			p.log.Print("context done, shutting down")
			break
		}
	}
}
