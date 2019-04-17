package flexMessager

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
)

func TestFlex(UserID string, bot *linebot.Client) {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "Hello,",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "World!",
				},
			},
		},
	}

	message := linebot.NewFlexMessage("alt text", container)
	_, err := bot.PushMessage(UserID, message).Do()
	if err != nil {
		log.Println("TestFlex: ", err)
	}
}
