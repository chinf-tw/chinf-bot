package messager

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

//EventTypeHandle 處理line-bot回應不同事件的對應狀況（例如：line-bot被加好友後的對應處理）
func EventTypeHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client, _temporaryStorage map[string][]string) {
	userid := event.Source.UserID
	switch event.Type {

	case linebot.EventTypeFollow:

		query := fmt.Sprintf("INSERT INTO spotify_user(line_id) VALUES ('%v') RETURNING id;", userid)
		dbQueryRow(db, query, userid, bot)
		PushMessage(userid, bot)

	case linebot.EventTypePostback:

		if event.Postback.Data == fmt.Sprintf("[%v][join member][yes]", userid) {
			// isRepeat := false
			// for _, TUserID := range _temporaryStorage["User_ID"] {
			// 	if TUserID == userid {
			// 		isRepeat = true
			// 		break
			// 	}
			// }
			// if !isRepeat {
			// 	_temporaryStorage["User_ID"] = append(_temporaryStorage["User_ID"], userid)
			// }
			query := `select build_JoinMember_cache($1::character varying(100));`
			// dbQueryRow(db, query, userid, bot)

			rows, err := db.Query(query, userid)
			if err != nil {
				log.Println("TpyeHandle 41 : ", err)
			}
			for rows.Next() {
				rows.Scan()
			}
			PushMessageSay(userid, bot, "請在“五分鐘內”輸入您的姓名，舉例：\n姓名為“王小明”\n就需輸入：[王小明]")
		}
	}
}

//MessageHandle 處理正在與使用者溝通的事件（例如：取得加入會員的名字）
func MessageHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client, _temporaryStorage map[string][]string) {
	var isPresence *bool
	query := fmt.Sprintf("SELECT is_presence('%v')", event.Source.UserID)
	rows, err := db.Query(query)
	if err != nil {
		log.Println("MessageHandle的err有問題找它！！--1", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&isPresence); err != nil {
			log.Println("MessageHandle的err有問題找它！！--2", err)
		}
	}
	if *isPresence {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			// var err error
			for _, TUserID := range _temporaryStorage["User_ID"] {
				println(TUserID)
				if TUserID == event.Source.UserID {
					query := fmt.Sprintf("UPDATE spotify_user SET name = '%v' WHERE line_id = '%v';", message.Text, event.Source.UserID)
					dbQueryRow(db, query, event.Source.UserID, bot)
				}
			}
			PushMessageSay(event.Source.UserID, bot, "恭喜成為會員！")
			// var reply string
			// if err == nil {
			// 	reply = "恭喜成為會員！"
			// } else {
			// 	reply = "出了一點問題，詢問一下工程師這發生什麼事吧。"
			// 	log.Println(err)
			// }

			// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			// 	log.Print(err)
			// }
		}
	}

}

func dbQueryRow(db *sql.DB, query string, userid string, bot *linebot.Client) (err error) {
	var response interface{}
	if err := db.QueryRow(query).Scan(&response); err != nil {
		sayErr := "出了一點問題，詢問一下工程師這發生什麼事吧。"
		log.Println(query, " ＜＝出問題！\n", err)
		PushMessageSay(userid, bot, sayErr)
		return err
	}
	// println(query)
	// println(userid)

	log.Printf("%v 對資料庫進行了 %v，資料庫回應為：%v", userid, query, response)
	return nil
}
