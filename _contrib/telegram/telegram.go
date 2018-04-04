package telegram

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoDialTelegram))   // connect the bot to telegram
	RegisterEvalGate(new(GoAcceptTelegram)) // accept a message
	RegisterEvalGate(new(GoReplyTelegram))  // reply to an accepted message
}

type GoDialTelegram struct {
	Token    string `ko:"name=token"` // telegram bot token
	Offset   int64  `ko:"name=offset"`
	PollNano int64  `ko:"name=pollNano"` // nanoseconds between long polls
}

func (dial GoDialTelegram) Play(ctx *runtime.Context) Conn {
	if bot, err := tb.NewBot(
		tb.Settings{
			Token: dial.Token,
			Poller: &tb.LongPoller{
				LastUpdateID: int(dial.Offset),
				Timeout:      time.Duration(dial.PollNano),
			},
		},
	); err != nil {
		log.Panicf("dial error (%v)", err)
		return nil
	} else {
		log.Printf("connected to telegram server as @%s", bot.Me.Username)
		return newConn(bot)
	}
}

type GoAcceptTelegram struct {
	On Conn `ko:"name=on,monadic"`
}

func (accept GoAcceptTelegram) Play(ctx *runtime.Context) Message {
	return *accept.On.Conn().Accept()
}

type GoReplyTelegram struct {
	To   Message `ko:"name=to"`
	With string  `ko:"name=with"`
}

func (reply GoReplyTelegram) Play(ctx *runtime.Context) struct{} {
	reply.To.On.Conn().Reply(reply.To, reply.With)
	return struct{}{}
}
