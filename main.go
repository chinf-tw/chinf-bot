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
	"chinf-bot/messager"
	"chinf-bot/userinfo"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// messager "github.com/chinf1996/Line-bot-messager"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

var botGlobal *linebot.Client
var temporaryStorage map[string][]string

func main() {

	temporaryStorage = map[string][]string{"User_ID": []string{}}
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	botGlobal = bot
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", lineCallbackHandler)
	http.HandleFunc("/chinf", selfCallbackHandler)
	http.HandleFunc("/test", testCallbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func lineCallbackHandler(w http.ResponseWriter, r *http.Request) {

	events, err := botGlobal.ParseRequest(r)

	judgeCallBackReq(w, err)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Println(err)
	}

	for _, event := range events {

		messager.EventTypeHandle(event, db, botGlobal, temporaryStorage)
		messager.MessageHandle(event, db, botGlobal, temporaryStorage)

		//(測試中)嘗試取得使用者給予之Line訊息
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Println(event.Source.UserID, "say : ", message.Text)
				if _, err = botGlobal.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

}

func selfCallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
	messager.JoinMember(db, botGlobal)
}

//testCallbackHandler 測試用區域
func testCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
	//測試加入會員功能用
	messager.PushMessage(os.Getenv("chinf_line_id"), botGlobal)
	userinfo.GetImage(botGlobal, db)
}
