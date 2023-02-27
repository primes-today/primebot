package primebot

import (
	"context"
	"log"
	"math/big"
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
	Logger         *log.Logger
	ServiceTimeout time.Duration
}

func NewPrimeBot(f Fetcher, t Ticker, g Generator, p Poster, opts *BotOpts) *PrimeBot {
	if opts == nil {
		opts = &BotOpts{}
	}

	if opts.Logger == nil {
		opts.Logger = log.New(&noopWriter{}, "", 0)
	}
	if opts.ServiceTimeout == 0 {
		opts.ServiceTimeout = 30 * time.Second
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

	fetchctx, cancel := context.WithTimeout(
		context.Background(),
		p.opts.ServiceTimeout,
	)
	cancel() // cancel to avoid leaking
	cur, err := p.ftc.Fetch(fetchctx)
	if err != nil {
		return err
	}

	next := cur.Posted.Add(p.tck.Interval())
	st := time.Until(next)
	p.log.Printf("retrieved status \"%d\"; next post in %v", cur.LastStatus, st)

	p.gen.SetStart(cur.LastStatus)
	t := p.tck.Start(ctx, st)
	pc := make(chan *big.Int, 1)
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
			postctx, cancel := context.WithTimeout(
				context.Background(),
				p.opts.ServiceTimeout,
			)
			err := p.pst.Post(postctx, status)
			cancel() // cancel to avoid leaking
			if err != nil {
				return err
			}
			p.log.Printf("posted status \"%d\"; next post in %v", status, p.tck.Interval())
		case err := <-er:
			p.log.Printf("error received from generator, %v. stopping", err)
			return err
		case <-ctx.Done():
			p.log.Print("context done, shutting down")
			return nil
		}
	}
}
