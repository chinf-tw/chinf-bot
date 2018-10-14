package messager

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

//EventTypeHandle 處理line-bot回應不同事件的對應狀況（例如：line-bot被加好友後的對應處理）
func EventTypeHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client, _temporaryStorage map[string][]string) {
	switch event.Type {

	case linebot.EventTypeFollow:

		query := fmt.Sprintf("INSERT INTO spotify_user( line_id,creation_time) VALUES ('%v',) RETURNING id;", event.Source.UserID)
		dbQueryRow(db, query)
		PushMessage(event.Source.UserID, bot)

	case linebot.EventTypePostback:

		if event.Postback.Data == fmt.Sprintf("[%v][yes]", event.Source.UserID) {
			isRepeat := false
			for _, TUserID := range _temporaryStorage["User_ID"] {
				if TUserID == event.Source.UserID {
					isRepeat = true
					break
				}
			}
			if !isRepeat {
				_temporaryStorage["User_ID"] = append(_temporaryStorage["User_ID"], event.Source.UserID)
			}
			PushMessageSay(event.Source.UserID, bot, "請輸入您的姓名")
		}
	}
}

//MessageHandle 處理正在與使用者溝通的事件（例如：取得加入會員的名字）
func MessageHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client, _temporaryStorage map[string][]string) {

	switch message := event.Message.(type) {
	case *linebot.TextMessage:

		for _, TUserID := range _temporaryStorage["User_ID"] {
			println(TUserID)
			if TUserID == event.Source.UserID {
				query := fmt.Sprintf("UPDATE spotify_user SET name = '%v' WHERE line_id = '%v';", message.Text, event.Source.UserID)
				dbQueryRow(db, query)
			}
		}

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
			log.Print(err)
		}
	}
}

func dbQueryRow(db *sql.DB, query string) {
	var userid int
	if err := db.QueryRow(query).Scan(&userid); err != nil {
		log.Println(query, " ＜＝出問題！\n", err)
	}
	println(query)
	println(userid)
}
