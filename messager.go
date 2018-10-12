package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

//PushMessage 輸出訊息用
func PushMessage(event *linebot.Event, bot *linebot.Client) {
	// _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("\\uDBC0\\uDC84 LINE emoji")).Do()
	// if err != nil {
	// 	print(err)
	// }
	// _, err = bot.PushMessage(event.Source.UserID, linebot.NewStickerMessage("1", "11")).Do()
	// if err != nil {
	// 	print(err)
	// }
	// yesBtn := linebot.NewMessageAction("left", "left clicked")
	yesBtn := linebot.NewMessageTemplateAction("yes", "我願意加入")
	butTemplate := linebot.NewButtonsTemplate("https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxbYZs9-LkRScXKWthdxw8gwUDUBkG34q0DgZnkI1pOkfybDx-",
		"來註冊成為我們的會員吧！",
		"有跟我一起買Spotify的夥伴就來成為會員吧", yesBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", butTemplate)
	_, err := bot.PushMessage(event.Source.UserID, message).Do()
	if err != nil {
		print(err)
	}
}
