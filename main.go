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
	"net/http"
	"os"
	"strings"
	"net/url"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.ToLower(message.Text) == "mid" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您的MID為："+event.Source.UserID)).Do(); err != nil {
						log.Print(err)
					}
				}
				var urlStr string = os.Getenv("ApiUrl")+"?MID="+event.Source.UserID+"&Text="+message.Text+"&GID="+event.Source.GroupID+"&RID="+event.Source.RoomID
				
				l3, err3 := url.Parse(urlStr)
				if err3 != nil {
					log.Fatal(err3)
				} else {
					var urlStr2 string = l3.Query().Encode()
					response, err := http.Get(os.Getenv("ApiUrl") + "?" + urlStr2)
					if err != nil {
						log.Fatal(err)
					} else {
						defer response.Body.Close()
						log.Print(response.Body)
					}
				}
			}
		}
	}
}
