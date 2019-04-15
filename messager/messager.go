package messager

import (
	"database/sql"
	"fmt"
	"log"

	"chinf-bot/userinfo"
	"github.com/line/line-bot-sdk-go/linebot"
)

// func init() {
// 	log.Println("Messager is camming...")
// }

//PushMessage 詢問是否要加入會員
func PushMessage(UserID string, bot *linebot.Client) {
	text := fmt.Sprintf("[%v][join member][yes]", UserID)
	yesBtn := linebot.NewPostbackAction("點我加入會員或更改姓名", text, "", "")
	butTemplate := linebot.NewButtonsTemplate(
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxbYZs9-LkRScXKWthdxw8gwUDUBkG34q0DgZnkI1pOkfybDx-",
		"來註冊成為我們的會員吧！",
		"有跟我一起買Spotify的夥伴就來成為會員吧",
		yesBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", butTemplate)
	_, err := bot.PushMessage(UserID, message).Do()
	if err != nil {
		print(err)
	}
}

//PushMessageSay 可以藉由此講出你想講的話
func PushMessageSay(UserID string, bot *linebot.Client, say string) {
	message := linebot.NewTextMessage(say)
	_, err := bot.PushMessage(UserID, message).Do()
	if err != nil {
		// print(err)
		log.Println(err)
	}
}

//JoinMember 處理已經有加好友，但還沒加會員的使用者
func JoinMember(db *sql.DB, bot *linebot.Client) {
	query := `select line_id from spotify_user where name is null;`

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Println(err)
		}
		PushMessage(userID, bot)
	}
}

func CarouselTemplate(UserID string, bot *linebot.Client, db *sql.DB) {
	var columns []*linebot.CarouselColumn

	userProfiles := userinfo.GetImages(bot, db)
	btn := linebot.NewMessageAction("test1", "test2")
	for _, userProfile := range userProfiles {
		column := linebot.NewCarouselColumn(userProfile.PictureURL, userProfile.DisplayName, userProfile.DisplayName, btn)
		columns = append(columns, column)
	}

	carouselTemplate := linebot.NewCarouselTemplate(columns...)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", carouselTemplate)
	_, err := bot.PushMessage(UserID, message).Do()
	if err != nil {
		log.Println("CarouselTemplate: ", err)
	}
}
