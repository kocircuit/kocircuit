package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoDialTwitter)) // connect to twitter
	RegisterEvalGate(new(GoHomeTimelineTwitter))
	RegisterEvalGate(new(GoUserTimelineTwitter))
}

type GoDialTwitter struct {
	ConsumerKey    string `ko:"name=consumerKey"`
	ConsumerSecret string `ko:"name=consumerSecret"`
	AccessToken    string `ko:"name=accessToken"`
	AccessSecret   string `ko:"name=accessSecret"`
}

func (dial GoDialTwitter) Play(ctx *runtime.Context) Conn {
	config := oauth1.NewConfig(dial.ConsumerKey, dial.ConsumerSecret)
	token := oauth1.NewToken(dial.AccessToken, dial.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token) // http.Client will automatically authorize Requests
	return newConn(twitter.NewClient(httpClient))
}

type GoHomeTimelineTwitter struct {
	Conn    Conn  `ko:"name=conn"`
	Count   int64 `ko:"name=count"`
	SinceID int64 `ko:"name=sinceID"`
	MaxID   int64 `ko:"name=maxID"`
}

type HomeTimelineResult struct {
	Tweets []*Tweet `ko:"name=tweets"`
	Error  *string  `ko:"name=error"`
}

func (g GoHomeTimelineTwitter) Play(ctx *runtime.Context) HomeTimelineResult {
	c := g.Conn.conn()
	if tweets, err := c.HomeTimeline(g.Count, g.SinceID, g.MaxID); err != nil {
		return HomeTimelineResult{Error: PtrString(err.Error())}
	} else {
		return HomeTimelineResult{Tweets: NewTweets(tweets)}
	}
}

type GoUserTimelineTwitter struct {
	Conn    Conn   `ko:"name=conn"`
	UserID  int64  `ko:"name=userID"`
	Screen  string `ko:"name=screen"`
	Count   int64  `ko:"name=count"`
	SinceID int64  `ko:"name=sinceID"`
	MaxID   int64  `ko:"name=maxID"`
}

type UserTimelineResult struct {
	Tweets []*Tweet `ko:"name=tweets"`
	Error  *string  `ko:"name=error"`
}

func (g GoUserTimelineTwitter) Play(ctx *runtime.Context) UserTimelineResult {
	c := g.Conn.conn()
	if tweets, err := c.UserTimeline(g.UserID, g.Screen, g.Count, g.SinceID, g.MaxID); err != nil {
		return UserTimelineResult{Error: PtrString(err.Error())}
	} else {
		return UserTimelineResult{Tweets: NewTweets(tweets)}
	}
}
