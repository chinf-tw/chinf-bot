// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	messager "github.com/chinf1996/Line-bot-messager"
	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

var botGlobal *linebot.Client
var selfevent *linebot.Event
var temporaryStorage map[string][]string

func main() {

	temporaryStorage = map[string][]string{"User_ID": []string{}}
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	botGlobal = bot
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/chinf", selfcallbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	events, err := botGlobal.ParseRequest(r)

	if err == nil {

		w.WriteHeader(200)

	} else {

		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Println(err)
	}

	for _, event := range events {

		selfevent = event

		if event.Type == linebot.EventTypeFollow {
			query := fmt.Sprintf("INSERT INTO spotify_user( line_id) VALUES ('%v') RETURNING id;", event.Source.UserID)
			var userid int
			err = db.QueryRow(query).Scan(&userid)
			if err != nil {
				log.Println(err)
			}
		} else if event.Type == linebot.EventTypePostback {
			if event.Postback.Data == fmt.Sprintf("[%v][yes]", event.Source.UserID) {
				isRepeat := false
				for _, TUserID := range temporaryStorage["User_ID"] {
					if TUserID == event.Source.UserID {
						isRepeat = true
						break
					}
				}
				if !isRepeat {
					temporaryStorage["User_ID"] = append(temporaryStorage["User_ID"], event.Source.UserID)
				}

				messager.PushMessageSay(event.Source.UserID, botGlobal, "請輸入您的姓名")
			}
		}
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			if _, err = botGlobal.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				log.Print(err)
			}

			for _, TUserID := range temporaryStorage["User_ID"] {
				println(TUserID)
				if TUserID == event.Source.UserID {
					var str string
					query := fmt.Sprintf("UPDATE spotify_user SET name = '%v' WHERE line_id = '%v';", message.Text, event.Source.UserID)
					if err := db.QueryRow(query).Scan(&str); err != nil {
						log.Println(err)
					}
					println(query)
					println(str)
				}
			}
		}
	}

}

func selfcallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
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
		messager.PushMessage(userID, botGlobal)
	}

}
