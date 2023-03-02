package primebot

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/mattn/go-mastodon"
)

type mockMastodon struct {
	getStatus chan []*mastodon.Status
	getId     chan mastodon.ID
	getError  chan error
	postError chan error
	postToot  chan *mastodon.Toot
}

func (m *mockMastodon) GetAccountStatuses(ctx context.Context, id mastodon.ID, pagination *mastodon.Pagination) ([]*mastodon.Status, error) {
	m.getId <- id
	select {
	case s := <-m.getStatus:
		return s, nil
	case e := <-m.getError:
		return nil, e
	}
}

func (m *mockMastodon) PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error) {
	m.postToot <- toot
	select {
	case e := <-m.postError:
		return nil, e
	default:
		return nil, nil
	}
}

func TestMastodonFetch(t *testing.T) {
	sc := make(chan []*mastodon.Status, 1)
	ic := make(chan mastodon.ID, 1)
	client := &mockMastodon{
		sc,
		ic,
		make(chan error),
		make(chan error),
		make(chan *mastodon.Toot),
	}
	m := &MastodonClient{
		client,
		"primes",
	}

	now := time.Now()
	statuses := []*mastodon.Status{
		{
			CreatedAt: now,
			Content:   "<p>1050167</p>",
		},
	}
	sc <- statuses

	status, err := m.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if status.LastStatus.Int64() != 1050167 {
		t.Errorf("got unexpected status: %s", status.LastStatus)
	}
	if status.Posted != now {
		t.Errorf("got unexpected post date: %s", status.Posted)
	}

	if id := <-ic; id != "primes" {
		t.Errorf("got unexpected id: %s", id)
	}
}

func TestMastodonFetchParseErr(t *testing.T) {
	sc := make(chan []*mastodon.Status, 1)
	ic := make(chan mastodon.ID, 1)
	client := &mockMastodon{
		sc,
		ic,
		make(chan error),
		make(chan error),
		make(chan *mastodon.Toot),
	}
	m := &MastodonClient{
		client,
		"primes",
	}

	now := time.Now()
	statuses := []*mastodon.Status{
		{
			CreatedAt: now,
			Content:   "1050167",
		},
	}
	sc <- statuses

	_, err := m.Fetch(context.Background())
	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestMastodonFetchErr(t *testing.T) {
	ec := make(chan error, 1)
	client := &mockMastodon{
		make(chan []*mastodon.Status, 1),
		make(chan mastodon.ID, 1),
		ec,
		make(chan error),
		make(chan *mastodon.Toot),
	}
	m := &MastodonClient{
		client,
		"primes",
	}

	ec <- errors.New("ack")

	_, err := m.Fetch(context.Background())
	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestMastodonPost(t *testing.T) {
	tc := make(chan *mastodon.Toot, 1)
	client := &mockMastodon{
		make(chan []*mastodon.Status),
		make(chan mastodon.ID),
		make(chan error),
		make(chan error),
		tc,
	}
	m := &MastodonClient{
		client,
		"primes",
	}

	err := m.Post(context.Background(), big.NewInt(3))
	if err != nil {
		t.Fatal(err)
	}

	s := <-tc
	if s.Status != "3" {
		t.Errorf("got unexpected status: %s", s.Status)
	}
}

func TestMastodonPostErr(t *testing.T) {
	ec := make(chan error, 1)
	client := &mockMastodon{
		make(chan []*mastodon.Status),
		make(chan mastodon.ID),
		make(chan error),
		ec,
		make(chan *mastodon.Toot, 1),
	}
	m := &MastodonClient{
		client,
		"primes",
	}

	ec <- errors.New("ugh")

	err := m.Post(context.Background(), big.NewInt(3))
	if err == nil {
		t.Error("expected error, got none")
	}
}
