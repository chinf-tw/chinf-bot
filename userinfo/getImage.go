package userinfo

import (
	"database/sql"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

//GetImages 取得會員的照片（目前為測試中，只會在後台看到）
func GetImages(bot *linebot.Client, db *sql.DB) (responses []*linebot.UserProfileResponse) {
	userIDList := getUserID(db)
	for _, userID := range userIDList {
		res, err := bot.GetProfile(userID).Do()
		if err != nil {
			log.Println(err)
		}
		// println(res.DisplayName)
		// println(res.PictureURL)
		// println(res.StatusMessage)
		responses = append(responses, res)

	}
	return responses
}
