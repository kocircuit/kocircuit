package twitter

import (
	"sync"

	"github.com/dghubble/go-twitter/twitter"
)

type Conn interface {
	conn() *conn
}

type conn struct {
	sync.Mutex
	Client *twitter.Client `ko:"name=client"`
}

func newConn(client *twitter.Client) *conn {
	return &conn{Client: client}
}

func (c *conn) conn() *conn {
	return c
}

func (c *conn) HomeTimeline(count, sinceID, maxID int64) ([]twitter.Tweet, error) {
	params := &twitter.HomeTimelineParams{
		Count: int(count), SinceID: sinceID, MaxID: maxID,
	}
	if tweets, _, err := c.Client.Timelines.HomeTimeline(params); err != nil {
		return nil, err
	} else {
		return tweets, nil
	}
}

func (c *conn) UserTimeline(userID int64, screen string, count, sinceID, maxID int64) ([]twitter.Tweet, error) {
	params := &twitter.UserTimelineParams{
		UserID: userID, ScreenName: screen,
		Count: int(count), SinceID: sinceID, MaxID: maxID,
	}
	if tweets, _, err := c.Client.Timelines.UserTimeline(params); err != nil {
		return nil, err
	} else {
		return tweets, nil
	}
}
