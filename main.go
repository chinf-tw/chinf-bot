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
	"fmt"
	"log"
	// "messager"
	"net/http"
	"os"

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

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		selfevent = event
		PushMessage(selfevent, botGlobal)
	}

}

func selfcallbackHandler(w http.ResponseWriter, r *http.Request) {
	PushMessage(selfevent, botGlobal)
}

// func before_main() {
// 	for _, event := range events {
// 		if event.Type == linebot.EventTypeMessage {
// 			switch message := event.Message.(type) {
// 			case *linebot.TextMessage:

// 				if message.Text == "menu" {
// 					leftBtn := linebot.NewMessageTemplateAction("left", "left clicked")
// 					rightBtn := linebot.NewMessageTemplateAction("right", "right clicked")
// 					template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
// 					userID := event.Source.UserID
// 					UserID_message := linebot.NewTextMessage(userID)

// 					messageeeee := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)

// 					_, MGerr := bot.PushMessage(userID, UserID_message).Do()

// 					if MGerr != nil {
// 						log.Println(MGerr)
// 					}

// 					_, MGerrr := bot.ReplyMessage(event.ReplyToken, messageeeee).Do()

// 					if MGerrr != nil {
// 						log.Println(MGerr)
// 					}
// 				}
// 				switch message.Text {
// 				case "left clicked":
// 					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("去你的 left clicked，閉嘴！")).Do()
// 				case "right clicked":
// 					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("去你的 right clicked，閉嘴！")).Do()
// 				}

// 				run(message.Text)

// 				switch runMassage {
// 				case 1:
// 					if !(isSay) {
// 						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("是否在説 "+name)).Do(); err != nil {
// 							log.Print(err)
// 						}
// 						isSay = true
// 					}
// 				case 2:
// 					outputMassage := ""
// 					switch porintIndexs {
// 					case 1:
// 						outputMassage = name + "你很棒"
// 					case 2:
// 						outputMassage = name + "你很帥的啦"
// 					case 3:
// 						outputMassage = name + "你很可愛喔"
// 					default:
// 					}
// 					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outputMassage)).Do(); err != nil {
// 						log.Print(err)
// 					}
// 					runMassage = 0
// 					checking = false
// 					isSay = false
// 					name = ""
// 					porintIndexs = 0
// 				default:

// 				}

// 			}
// 		}
// 	}
// }

// func checked(str string) {
// 	if str == "left clicked" {

// 	}
// }

// func run(str string) {
// 	porintName, isHave := massage_bot(str)
// 	if isHave {
// 		checking = true
// 		name = porintName
// 		runMassage = 1
// 	}

// 	checkmassage(str)
// }
// func checkmassage(strr string) {
// 	var yes = []string{"yes", "是", "沒錯", "當然啊", "對啦", "對啊", "當然啦", "對"}
// 	_, x := porintSearch(strr, yes)
// 	if checking && x {
// 		runMassage = 2
// 	}

// }
// func massage_bot(massage string) (string, bool) {
// 	var porint = [][]string{{"志嘉", "思毅", "淂堉"}, {"洪甫", "權甫"}, {"敏文", "文文", "阿文"}}

// 	return porintSearch2(massage, porint)
// }
// func porintSearch(massage string, porint []string) (string, bool) {
// 	var inside = ""
// 	var ishavePorint = false
// 	for _, element := range porint {
// 		if strings.Contains(massage, element) {
// 			inside += element + ""
// 			ishavePorint = true

// 		}
// 	}
// 	return inside, ishavePorint
// }
// func porintSearch2(massage string, porint [][]string) (string, bool) {
// 	var inside = ""
// 	var ishavePorint = false

// 	for index, element := range porint {
// 		for _, ele := range element {
// 			if strings.Contains(massage, ele) {
// 				inside += ele + ""
// 				ishavePorint = true
// 				porintIndexs = index + 1
// 			}
// 		}

// 	}
// 	return inside, ishavePorint
// }
