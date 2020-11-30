package main

import (
	"os"

	"github.com/go-chat-bot/bot/slack"
	"github.com/go-chat-bot/bot/telegram"
	_ "github.com/hcaijin/sgk-bot/plugin"
	_ "github.com/go-chat-bot/plugins/uptime"
)

func main() {
    go telegram.Run(os.Getenv("TG_TOKEN"), os.Getenv("DEBUG") != "")
    slack.Run(os.Getenv("SLACK_TOKEN"))
}
