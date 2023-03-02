package primebot

import (
	"context"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type mockTwitter struct {
	fetchResponse chan []twitter.Tweet
	userParam     chan *twitter.UserTimelineParams
	statusParam   chan string
	fetchError    chan error
	updateError   chan error
}

func (t *mockTwitter) UserTimeline(user *twitter.UserTimelineParams) ([]twitter.Tweet, *http.Response, error) {
	t.userParam <- user
	select {
	case r := <-t.fetchResponse:
		return r, nil, nil
	case e := <-t.fetchError:
		return nil, nil, e
	}
}

func (t *mockTwitter) Update(status string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error) {
	t.statusParam <- status
	select {
	case e := <-t.updateError:
		return nil, nil, e
	default:
		return nil, nil, nil
	}
}

func TestTwitterFetch(t *testing.T) {
	dtStr := time.Now().Format(time.RubyDate)
	dt, _ := time.Parse(time.RubyDate, dtStr)
	res := []twitter.Tweet{
		{
			Text:      "1050139",
			CreatedAt: dtStr,
			FullText:  "1050139",
		},
	}
	tc := make(chan []twitter.Tweet, 1)
	up := make(chan *twitter.UserTimelineParams, 1)
	client := &mockTwitter{
		tc,
		up,
		make(chan string),
		make(chan error),
		make(chan error),
	}

	tc <- res

	user := &twitter.User{
		ID:    1000,
		IDStr: "1000",
	}
	fetcher := TwitterClient{
		t: client,
		s: client,
		u: user,
	}

	status, err := fetcher.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if s := status.LastStatus.Int64(); s != 1050139 {
		t.Errorf("got unexpected status %d", s)
	}
	if status.Posted != dt {
		t.Errorf("got unexpected date %s", status.Posted)
	}

	params := <-up
	if user.ID != params.UserID {
		t.Errorf("unexpected user id %d", params.UserID)
	}
}

func TestTwitterPost(t *testing.T) {
	s := make(chan string, 1)
	client := &mockTwitter{
		make(chan []twitter.Tweet),
		make(chan *twitter.UserTimelineParams, 1),
		s,
		make(chan error),
		make(chan error),
	}
	user := &twitter.User{
		ID:    1000,
		IDStr: "1000",
	}

	status := big.NewInt(1050139)
	poster := TwitterClient{
		t: client,
		s: client,
		u: user,
	}
	err := poster.Post(context.Background(), status)
	if err != nil {
		t.Fatal(err)
	}

	posted := <-s
	if posted != "1050139" {
		t.Errorf("unexpected status posted: %s", posted)
	}
}
