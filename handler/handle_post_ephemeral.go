package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type EphemeralRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

// 特定のユーザーにだけ見えるメッセージを送信する
func HandlePostEphemeral(api *slack.Client, channelID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req EphemeralRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.UserID == "" || req.Text == "" {
			http.Error(w, "Invalid Request (user_id と text が必要です)", http.StatusBadRequest)
			return
		}

		// Slack API: chat.postEphemeral
		// 送り先のチャンネルID、本人しか見えないようにするためのUserID、テキストを指定する
		_, err = api.PostEphemeral(channelID, req.UserID, slack.MsgOptionText(req.Text, false))
		if err != nil {
			log.Printf("chat.postEphemeral エラー: %v", err)
			http.Error(w, "Ephemeralメッセージの送信に失敗しました", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}
}
