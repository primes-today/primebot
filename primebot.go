package primebot

import (
	"context"
	"time"
)

type Bot interface {
	Start(context.Context) error
}

type BotOpts struct{}

func NewPrimeBot(f Fetcher, t Ticker, g Generator, r Formatter, p Poster, opts *BotOpts) *PrimeBot {
	if opts == nil {
		opts = &BotOpts{}
	}
	return &PrimeBot{
		ftc:  f,
		tck:  t,
		gen:  g,
		fmt:  r,
		pst:  p,
		opts: opts,
	}
}

type PrimeBot struct {
	ftc  Fetcher
	tck  Ticker
	gen  Generator
	fmt  Formatter
	pst  Poster
	opts *BotOpts
}

func (p *PrimeBot) Start(ctx context.Context) error {
	cur, err := p.ftc.Fetch()
	if err != nil {
		return err
	}

	st := cur.ts.Sub(time.Now())

	p.gen.SetStart(cur.num)
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
			err := p.pst.Post(p.fmt.Format(<-pc))
			if err != nil {
				return err
			}
		case err := <-er:
			return err
		case <-ctx.Done():
			break
		}
	}
}
