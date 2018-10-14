package main

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func judgeCallBackReq(w http.ResponseWriter, err error) {
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
}
