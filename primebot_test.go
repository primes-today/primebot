package primebot

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"
)

type MockFetcher struct {
	s chan *Status
	e chan error
}

func (f *MockFetcher) Fetch(ctx context.Context) (*Status, error) {
	select {
	case s := <-f.s:
		return s, nil
	case e := <-f.e:
		return nil, e
	}
}

type MockTicker struct {
	t chan time.Time
	d chan time.Duration
	i time.Duration
}

func (t *MockTicker) Start(ctx context.Context, first time.Duration) (c chan time.Time) {
	t.d <- first
	return t.t
}

func (t *MockTicker) Interval() time.Duration {
	return t.i
}

type MockGenerator struct {
	i chan *big.Int
	e chan error
}

func (g *MockGenerator) SetStart(i *big.Int) {}

func (g *MockGenerator) Generate(ctx context.Context) (*big.Int, error) {
	select {
	case i := <-g.i:
		return i, nil
	case e := <-g.e:
		return nil, e
	}
}

type MockPoster struct {
	i chan *big.Int
	e chan error
}

func (p *MockPoster) Post(ctx context.Context, status *big.Int) error {
	select {
	case e := <-p.e:
		return e
	default:
		p.i <- status
		return nil
	}
}

func TestPrimebot(t *testing.T) {
	fc := make(chan *Status, 1)
	fetcher := &MockFetcher{
		fc,
		make(chan error),
	}
	tc := make(chan time.Time, 1)
	dc := make(chan time.Duration, 1)
	ticker := &MockTicker{
		tc,
		dc,
		time.Second * 10,
	}
	gc := make(chan *big.Int, 1)
	generator := &MockGenerator{
		gc,
		make(chan error),
	}
	pc := make(chan *big.Int, 1)
	poster := &MockPoster{
		pc,
		make(chan error),
	}
	bot := NewPrimeBot(fetcher, ticker, generator, poster, nil)

	fc <- &Status{
		LastStatus: big.NewInt(5),
		Posted:     time.Now(),
	}
	tc <- time.Now()
	gc <- big.NewInt(7)

	ctx, cancel := context.WithCancel(context.Background())
	err := make(chan error)
	done := make(chan interface{})
	go func() {
		if e := bot.Start(ctx); e != nil {
			err <- e
		}
		done <- nil
	}()

	posted := <-pc
	if posted.Cmp(big.NewInt(7)) != 0 {
		t.Errorf("expected 7 to be posted, got %s", posted)
	}

	select {
	case <-pc:
		t.Fatal("got new post, when none should've been posted yet")
	default:
	}

	tc <- time.Now()
	gc <- big.NewInt(11)
	posted = <-pc
	if posted.Cmp(big.NewInt(11)) != 0 {
		t.Errorf("expected 11 to be posted, got %s", posted)
	}

	cancel()

	select {
	case <-done:
		break
	case e := <-err:
		t.Fatal(e)
	}
}

func TestPrimebotFetchFailure(t *testing.T) {
	fe := make(chan error, 1)
	fetcher := &MockFetcher{
		make(chan *Status),
		fe,
	}
	ticker := &MockTicker{
		make(chan time.Time),
		make(chan time.Duration),
		time.Second * 10,
	}
	generator := &MockGenerator{
		make(chan *big.Int),
		make(chan error),
	}
	pc := make(chan *big.Int, 1)
	poster := &MockPoster{
		pc,
		make(chan error),
	}
	bot := NewPrimeBot(fetcher, ticker, generator, poster, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := make(chan error)
	go func() {
		if e := bot.Start(ctx); e != nil {
			err <- e
		}
	}()

	fe <- errors.New("frig")

	e := <-err
	if e.Error() != "frig" {
		t.Errorf("got unexpected error message %s", e)
	}

	select {
	case <-pc:
		t.Fatal("got post, when none should've been posted")
	default:
	}

}

func TestPrimebotGeneratorFail(t *testing.T) {
	fc := make(chan *Status, 1)
	fetcher := &MockFetcher{
		fc,
		make(chan error, 1),
	}
	ticker := &MockTicker{
		make(chan time.Time, 1),
		make(chan time.Duration, 1),
		time.Second * 10,
	}
	ge := make(chan error, 1)
	generator := &MockGenerator{
		make(chan *big.Int),
		ge,
	}
	pc := make(chan *big.Int, 1)
	poster := &MockPoster{
		pc,
		make(chan error, 1),
	}
	bot := NewPrimeBot(fetcher, ticker, generator, poster, nil)

	fc <- &Status{
		LastStatus: big.NewInt(5),
		Posted:     time.Now(),
	}
	ge <- errors.New("frig")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := make(chan error)
	go func() {
		if e := bot.Start(ctx); e != nil {
			err <- e
		}
	}()

	e := <-err
	if e.Error() != "frig" {
		t.Errorf("got unexpected error message %s", e)
	}

	select {
	case <-pc:
		t.Fatal("got post, when none should've been posted")
	default:
	}
}

func TestPrimebotPostFail(t *testing.T) {
	fc := make(chan *Status, 1)
	fetcher := &MockFetcher{
		fc,
		make(chan error),
	}
	tc := make(chan time.Time, 1)
	dc := make(chan time.Duration, 1)
	ticker := &MockTicker{
		tc,
		dc,
		time.Second * 10,
	}
	gc := make(chan *big.Int, 1)
	generator := &MockGenerator{
		gc,
		make(chan error),
	}
	pe := make(chan error, 1)
	poster := &MockPoster{
		make(chan *big.Int, 1),
		pe,
	}
	bot := NewPrimeBot(fetcher, ticker, generator, poster, nil)

	fc <- &Status{
		LastStatus: big.NewInt(5),
		Posted:     time.Now(),
	}
	tc <- time.Now()
	gc <- big.NewInt(7)
	pe <- errors.New("frig")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := make(chan error)
	go func() {
		if e := bot.Start(ctx); e != nil {
			err <- e
		}
	}()

	e := <-err
	if e.Error() != "frig" {
		t.Errorf("got unexpected error message %s", e)
	}
}
