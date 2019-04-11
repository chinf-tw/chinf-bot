package userinfo

import (
	"database/sql"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)
//Getimage 取得會員的照片（目前為測試中，只會在後台看到）
func GetImage(bot *linebot.Client, db *sql.DB) {
	userIDList := getUserID(db)
	for _, userID := range userIDList {
		res, err := bot.GetProfile(userID).Do()
		if err != nil {
			log.Println(err)
		}
		println(res.DisplayName)
		println(res.PictureURL)
		println(res.StatusMessage)
	}
}
