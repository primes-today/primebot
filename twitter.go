package primebot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

func NewTwitterClient(c *http.Client) (*TwitterClient, error) {
	client := twitter.NewClient(c)
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return &TwitterClient{
		c: client,
		u: user,
	}, nil
}

type TwitterClient struct {
	c *twitter.Client
	u *twitter.User
}

func (t *TwitterClient) Fetch(ctx context.Context) (*Status, error) {
	ss, _, err := t.c.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          t.u.ID,
		Count:           1,
		ExcludeReplies:  twitter.Bool(true),
		IncludeRetweets: twitter.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	if len(ss) < 1 {
		return nil, errors.New("unable to retrieve last posting")
	}

	s := ss[0]
	n, err := strconv.ParseUint(s.Text, 10, 64)
	if err != nil {
		return nil, err
	}

	ts, err := s.CreatedAtTime()
	if err != nil {
		return nil, err
	}

	return &Status{
		num: n,
		ts:  ts,
	}, nil
}

func (t *TwitterClient) Post(ctx context.Context, status uint64) error {
	_, _, err := t.c.Statuses.Update(fmt.Sprintf("%d", status), nil)
	return err
}
