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

// var checking bool = false
// var runMassage = 0
// var name = ""
// var isSay = false
// var porintIndexs = 0

func main() {

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
		}

		// println("準備進入messager階段")
		// messager.PushMessage(selfevent, botGlobal)
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
