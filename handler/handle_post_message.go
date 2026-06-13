package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type MessageRequest struct {
	Text string `json:"text"`
}

// メッセージ送信を処理する
func HandlePostMessage(api *slack.Client, channelID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req MessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Text == "" {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		// Slack API: chat.postMessage
		_, _, err = api.PostMessage(channelID, slack.MsgOptionText(req.Text, false))
		if err != nil {
			log.Printf("Slack送信エラー: %v", err)
			http.Error(w, "Slackへの送信に失敗しました", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}
}
