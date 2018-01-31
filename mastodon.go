package primebot

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	mastodon "github.com/mattn/go-mastodon"
)

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
	c  *mastodon.Client
	id mastodon.ID
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
	n, err := strconv.ParseUint(s.Content, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Status{
		num: n,
		ts:  s.CreatedAt,
	}, nil
}

func (m *MastodonClient) Post(ctx context.Context, status uint64) error {
	// _, err := m.c.PostStatus(ctx, &mastodon.Toot{
	// 	Status: fmt.Sprintf("%d", status),
	// })
	fmt.Printf("would post %d\n", status)

	// return err
	return nil
}
