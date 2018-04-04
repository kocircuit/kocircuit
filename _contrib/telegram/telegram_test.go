package telegram

import (
	"fmt"
	"os"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestTelegram(t *testing.T) {
	c := DialGate{
		Token:    os.Getenv("TELEGRAM_TOKEN"), // declare -x TELEGRAM_TOKEN=YourTelegramBotToken
		PollNano: 1e9,                         // 1 sec
	}.Play(nil)
	for {
		accepted := AcceptGate{On: c}.Play(nil)
		fmt.Println(Sprint(accepted))
		ReplyGate{To: accepted, With: fmt.Sprintf("ok, %s", accepted.Text)}.Play(nil)
	}
}
