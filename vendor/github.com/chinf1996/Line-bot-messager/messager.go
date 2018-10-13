package messager

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func init() {
	log.Println("Messager is camming...")
}

//PushMessage 輸出訊息用
func PushMessage(UserId string, bot *linebot.Client) {

	yesBtn := linebot.NewMessageAction("我願意", "[yes]")
	butTemplate := linebot.NewButtonsTemplate(
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxbYZs9-LkRScXKWthdxw8gwUDUBkG34q0DgZnkI1pOkfybDx-",
		"來註冊成為我們的會員吧！",
		"有跟我一起買Spotify的夥伴就來成為會員吧",
		yesBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", butTemplate)
	_, err := bot.PushMessage(UserId, message).Do()
	if err != nil {
		print(err)
	}
}
