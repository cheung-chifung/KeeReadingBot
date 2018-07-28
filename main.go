package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	botChannel := os.Getenv("BOT_CHANNEL")

	srv := NewBotServer(botToken, botChannel)
	if err := srv.Start(); err != nil {
		logrus.Fatal(err)
	}
}
