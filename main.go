// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// +build !appengine

package main

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

//func main() {
//
//	fmt.Println("YOYO")
//
//	channelSecret := "cf4558ccb428b620d350958982aff369"
//	channelAccessToken := "Obb9/VQi9jJSsLTFmnK0tKADZhM6vnDUa0qCEdK5G1t4e3bgszrWQMnGcrKXF1GTJTrdd92LijNTUq8sA1cP6IjrwZNivjWKkRpH/623CO5yPENHcX74i9oe0gkj6lPyDoCLnIEhLk41JtKQJwCDnwdB04t89/1O/w1cDnyilFU="
//	handler, err := httphandler.New(
//		channelSecret,
//		channelAccessToken,
//	)
//	if err != nil {
//		fmt.Println(err)
//		log.Fatal(err)
//	}
//
//	// Setup HTTP Server for receiving requests from LINE platform
//	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
//		bot, err := handler.NewClient()
//		if err != nil {
//			fmt.Println("aaa")
//			log.Print(err)
//			return
//		}
//		fmt.Println("bbb")
//		for _, event := range events {
//			fmt.Println("eee")
//			if event.Type == linebot.EventTypeMessage {
//				fmt.Println("ccc")
//				switch message := event.Message.(type) {
//				case *linebot.TextMessage:
//					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
//						fmt.Println(err)
//						log.Print(err)
//					} else {
//						fmt.Println("111", message)
//					}
//				default : fmt.Println("222", message)
//				}
//			} else {
//				fmt.Println("ddd")
//			}
//		}
//		fmt.Println("fff")
//	})
//	fmt.Println("ggg")
//	http.Handle("/callback", handler)
//	fmt.Println("hhh")
//	// This is just a sample code.
//	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
//	if err := http.ListenAndServe(":8877", nil); err != nil {
//		fmt.Println("iii")
//		log.Fatal(err)
//	}
//	fmt.Println("jjj")
//}

func main() {

	channelSecret := "cf4558ccb428b620d350958982aff369"
	channelAccessToken := "Obb9/VQi9jJSsLTFmnK0tKADZhM6vnDUa0qCEdK5G1t4e3bgszrWQMnGcrKXF1GTJTrdd92LijNTUq8sA1cP6IjrwZNivjWKkRpH/623CO5yPENHcX74i9oe0gkj6lPyDoCLnIEhLk41JtKQJwCDnwdB04t89/1O/w1cDnyilFU="

	handler, err := httphandler.New(
		channelSecret,
		channelAccessToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	http.Handle("/callback", handler)
	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal(err)
	}
}
