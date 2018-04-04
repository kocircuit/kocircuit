package twitter

import (
	"fmt"
	"os"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestTwitter(t *testing.T) {
	result := HomeTimelineGate{
		Conn: DialGate{
			ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
			ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
			AccessToken:    os.Getenv("TWITTER_ACCESS_TOKEN"),
			AccessSecret:   os.Getenv("TWITTER_ACCESS_SECRET"),
		}.Play(nil),
		Count:   3,
		SinceID: 0,
		MaxID:   0,
	}.Play(nil)
	fmt.Println(Sprint(result))
}
