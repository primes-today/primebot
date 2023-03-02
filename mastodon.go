package primebot

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	mastodon "github.com/mattn/go-mastodon"
)

var (
	MastodonStatusRex = regexp.MustCompile(`^<p>([0-9]+)</p>$`)
)

type MastodonPrimesClient interface {
	GetAccountStatuses(context.Context, mastodon.ID, *mastodon.Pagination) ([]*mastodon.Status, error)
	PostStatus(context.Context, *mastodon.Toot) (*mastodon.Status, error)
}

type MastodonConfig struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
	AccountID    string
}

func NewMastodonClient(ctx context.Context, config *MastodonConfig) (*MastodonClient, error) {
	mconfig := &mastodon.Config{
		Server:       config.Server,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		AccessToken:  config.AccessToken,
	}
	c := mastodon.NewClient(mconfig)

	act, err := c.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	return &MastodonClient{
		c:  c,
		id: act.ID,
	}, nil
}

type MastodonClient struct {
	c  MastodonPrimesClient
	id mastodon.ID
}

func (m *MastodonClient) parseStatus(status string) (*big.Int, error) {
	n := &big.Int{}

	ss := MastodonStatusRex.FindStringSubmatch(status)
	if l := len(ss); l != 2 {
		return n, fmt.Errorf("unexpected number of matches: %v", l)
	}
	s := ss[1]
	if s == "" {
		return n, errors.New("did not find substring match")
	}
	n, success := n.SetString(s, 10)
	if !success {
		return n, fmt.Errorf("unable to convert status to bigint: %s", s)
	}

	return n, nil
}

func (m *MastodonClient) Fetch(ctx context.Context) (*Status, error) {
	ss, err := m.c.GetAccountStatuses(ctx, m.id, nil)
	if err != nil {
		return nil, err
	}

	if len(ss) < 1 {
		return nil, errors.New("retrieved empty list of statuses")
	}

	s := ss[0]
	n, err := m.parseStatus(s.Content)
	if err != nil {
		return nil, err
	}

	return &Status{
		LastStatus: n,
		Posted:     s.CreatedAt,
	}, nil
}

func (m *MastodonClient) Post(ctx context.Context, status *big.Int) error {
	_, err := m.c.PostStatus(ctx, &mastodon.Toot{
		Status: status.Text(10),
	})

	return err
}
