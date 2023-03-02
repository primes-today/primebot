package primebot

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

type TwitterPrimesTimelineService interface {
	UserTimeline(*twitter.UserTimelineParams) ([]twitter.Tweet, *http.Response, error)
}

type TwitterPrimesStatusService interface {
	Update(string, *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error)
}

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
		t: client.Timelines,
		s: client.Statuses,
		u: user,
	}, nil
}

type TwitterClient struct {
	t TwitterPrimesTimelineService
	s TwitterPrimesStatusService
	u *twitter.User
}

func (t *TwitterClient) Fetch(_ context.Context) (*Status, error) {
	n := &big.Int{}
	ss, _, err := t.t.UserTimeline(&twitter.UserTimelineParams{
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
	n, success := n.SetString(s.Text, 10)
	if !success {
		return nil, fmt.Errorf("unable to convert status to bigint: %s", s.Text)
	}

	ts, err := s.CreatedAtTime()
	if err != nil {
		return nil, err
	}

	return &Status{
		LastStatus: n,
		Posted:     ts,
	}, nil
}

func (t *TwitterClient) Post(_ context.Context, status *big.Int) error {
	_, _, err := t.s.Update(status.Text(10), nil)
	return err
}
