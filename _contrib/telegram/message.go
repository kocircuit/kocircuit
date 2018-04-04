package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewMessage(on Conn, msg *tb.Message) *Message {
	if msg == nil {
		return nil
	} else {
		return &Message{
			On:           on,
			ID:           int64(msg.ID),
			Sender:       NewUser(msg.Sender),
			Chat:         NewChat(msg.Chat),
			OriginSender: NewUser(msg.OriginalSender),
			OriginChat:   NewChat(msg.OriginalChat),
			ReplyTo:      NewMessage(on, msg.ReplyTo),
			Signature:    msg.Signature,
			Text:         msg.Text,
			Payload:      msg.Payload,
		}
	}
}

type Message struct {
	On           Conn     `ko:"name=on"`
	ID           int64    `ko:"name=id"`
	Sender       *User    `ko:"name=sender"` // for messages sent to channels, Sender will be nil
	Chat         *Chat    `ko:"name=chat"`
	OriginSender *User    `ko:"name=originSender"` // for forwarded messages
	OriginChat   *Chat    `ko:"name=originChat"`
	ReplyTo      *Message `ko:"name=replyTo"`
	Signature    string   `ko:"name=signature"` // author signature in channels
	Text         string   `ko:"name=text"`
	Payload      string   `ko:"name=payload"` // command message payload
	// TODO: files, images, videos, contact, location, etc
}

func NewUser(user *tb.User) *User {
	if user == nil {
		return nil
	} else {
		return &User{
			ID:    int64(user.ID),
			First: user.FirstName,
			Last:  user.LastName,
			User:  user.Username,
		}
	}
}

type User struct {
	ID    int64  `ko:"name=id"`
	First string `ko:"name=first"`
	Last  string `ko:"name=last"`
	User  string `ko:"name=user"`
}

func (u *User) Translate() *tb.User {
	return &tb.User{
		ID:        int(u.ID),
		FirstName: u.First,
		LastName:  u.Last,
		Username:  u.User,
	}
}

func NewChat(chat *tb.Chat) *Chat {
	if chat == nil {
		return nil
	} else {
		return &Chat{
			ID:    chat.ID,
			Type:  string(chat.Type),
			Title: chat.Title,
			First: chat.FirstName,
			Last:  chat.LastName,
			User:  chat.Username,
		}
	}
}

type Chat struct {
	ID    int64  `ko:"name=id"`
	Type  string `ko:"name=type"`  // private, group, supergroup, channel, privatechannel
	Title string `ko:"name=title"` // empty for private chat
	First string `ko:"name=first"`
	Last  string `ko:"name=last"`
	User  string `ko:"name=user"`
}
