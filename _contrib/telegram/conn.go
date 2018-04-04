package telegram

import (
	"log"
	"sync"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Conn interface {
	Conn() *conn
}

type conn struct {
	sync.Mutex
	Bot    *tb.Bot          `ko:"name=bot"`
	Update chan *tb.Message `ko:"name=update"`
}

func newConn(bot *tb.Bot) *conn {
	up := make(chan *tb.Message)
	bot.Handle(tb.OnText, func(msg *tb.Message) { up <- msg })
	go func() { bot.Start() }()
	return &conn{Bot: bot, Update: up}
}

func (c *conn) Conn() *conn { return c }

func (c *conn) Accept() *Message {
	return NewMessage(c, <-c.Update)
}

func (c *conn) Reply(to Message, with string) {
	c.Lock()
	defer c.Unlock()
	if _, err := c.Bot.Send(to.Sender.Translate(), with); err != nil {
		log.Printf("send error (%v)", err)
	}
}
