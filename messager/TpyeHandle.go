package messager

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

//EventTypeHandle 處理line-bot回應不同事件的對應狀況（例如：line-bot被加好友後的對應處理）
func EventTypeHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client) {
	userid := event.Source.UserID
	switch event.Type {

	case linebot.EventTypeFollow:

		query := fmt.Sprintf("INSERT INTO spotify_user(line_id) VALUES ('%v') RETURNING id;", userid)
		dbQueryRow(db, query, userid, bot)
		PushMessage(userid, bot)

	case linebot.EventTypePostback:

		if event.Postback.Data == fmt.Sprintf("[%v][join member][yes]", userid) {

			query := `select build_JoinMember_cache($1::character varying(100));`
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
func MessageHandle(event *linebot.Event, db *sql.DB, bot *linebot.Client) {
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
			var userMessage = []rune(message.Text)
			//小解釋：預設得到的資料為[<name>]故需要先判斷開頭跟結尾是“[”跟“]”
			isLeft := (string(userMessage[0]) == "[" || string(userMessage[0]) == "［")
			isRight := (string(userMessage[len(userMessage)-1]) == "]" || string(userMessage[len(userMessage)-1]) == "］")
			if isLeft && isRight {
				//userMessage[1:len(userMessage)-1]的意思是去除掉[]的字詞
				query := fmt.Sprintf("UPDATE spotify_user SET name = '%v' WHERE line_id = '%v';", string(userMessage[1:len(userMessage)-1]), event.Source.UserID)
				err := dbQueryRow(db, query, event.Source.UserID, bot)
				if err == nil {
					PushMessageSay(event.Source.UserID, bot, "恭喜成為會員！/更改姓名成功！")
				}
			}
		}
	}

}

func dbQueryRow(db *sql.DB, query string, userid string, bot *linebot.Client) (err error) {
	var response interface{}
	err = db.QueryRow(query).Scan(&response)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("%v 對資料庫進行了 %v，且沒有回傳資料", userid, query)
	case err != nil:
		sayErr := "出了一點問題，詢問一下工程師這發生什麼事吧。"
		log.Println(query, " ＜＝出問題！\n", err)
		PushMessageSay(userid, bot, sayErr)
		return err
	default:
		log.Printf("%v 對資料庫進行了 %v，資料庫回應為：%v", userid, query, response)
	}
	return nil
}
