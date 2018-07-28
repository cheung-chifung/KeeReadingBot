package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotServer struct {
	botToken   string
	botChannel string
	bot        *tb.Bot
	cmtMgr     CommentManager
}

func NewBotServer(botToken string, botChannel string) *BotServer {
	return &BotServer{
		botToken:   botToken,
		botChannel: botChannel,
		cmtMgr:     NewCommentManager(),
	}
}

func (b *BotServer) Start() error {
	if b.bot != nil {
		return errors.New("server is started")
	}

	var err error
	b.bot, err = tb.NewBot(tb.Settings{
		Token:  b.botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}

	b.bot.Handle("/hi", b.HandleHi)
	b.bot.Handle("/start", b.HandleStart)
	b.bot.Handle(tb.OnText, b.HandleDefaultText)

	logrus.Info("Bot is ready")
	b.bot.Start()
	return nil
}

func (b *BotServer) HandleHi(msg *tb.Message) {
	b.bot.Send(msg.Sender, fmt.Sprintf("Hello, %s", msg.Sender.Username))
}

func (b *BotServer) HandleStart(msg *tb.Message) {
	b.cmtMgr.Reset()
	b.bot.Send(msg.Sender, "Please paste a new url")
}

func (b *BotServer) HandleDefaultText(msg *tb.Message) {
	switch b.cmtMgr.State() {
	case CommentStateWaiting:
		txt := msg.Text
		url, err := url.Parse(txt)
		if err != nil {
			logrus.WithError(err).Warn(fmt.Sprintf("wrong url: %s", txt))
			return
		}
		b.cmtMgr.InputURL(url)
		b.bot.Send(msg.Sender, "Please enter your comment")
	case CommentStateURLEntered:
		comment := msg.Text
		b.cmtMgr.InputComment(comment)
		b.bot.Send(msg.Sender, "Please confirm your comment.")
		b.bot.Send(msg.Sender, fmt.Sprintf("Link:    %s", b.cmtMgr.URL()))
		b.bot.Send(msg.Sender, fmt.Sprintf("Comment: %s", b.cmtMgr.Comment()))
		b.bot.Send(msg.Sender, "Are you sure?")
	case CommentStateCommentEntered:
		resp := strings.ToUpper(msg.Text)
		if resp == "YES" {
			err := b.PublishComment()
			if err != nil {
				b.bot.Send(msg.Sender, "Publish fail, enter 'Yes' to retry.")
				return
			}
			b.bot.Send(msg.Sender, "Your comment is published")
		} else if resp == "NO" {
			b.cmtMgr.Reset()
			b.bot.Send(msg.Sender, "Your draft is deleted")
		}
	}
}

func (b *BotServer) PublishComment() error {
	logrus.Printf("Publishing comment to %s", b.botChannel)
	ch := &tb.Chat{Username: b.botChannel, Type: tb.ChatChannel}
	b.bot.Send(ch, fmt.Sprintf("%s\n\n%s", b.cmtMgr.URL(), b.cmtMgr.Comment()))
	return nil
}

func (b *BotServer) Handle(msg *tb.Message) {
	b.bot.Send(msg.Sender, "Please paste and enter your comments")
}
