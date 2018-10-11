package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

//PushMessage 輸出訊息用
func PushMessage(event *linebot.Event, bot *linebot.Client) {
	_, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("TEST")).Do()
	if err != nil {
		print(err)
	}
}