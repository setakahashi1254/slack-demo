package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// HandlePostMessage - メッセージ送信と同時にスタンプ（リアクション）も自動で押す
func HandlePostMessageWithStamp(api *slack.Client, channelID string) http.HandlerFunc {
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

		// メッセージを送信
		// 返り値の「respTimestamp」が、このメッセージの固有の住所（ID）になる
		_, respTimestamp, err := api.PostMessage(channelID, slack.MsgOptionText(req.Text, false))
		if err != nil {
			log.Printf("Slack送信エラー: %v", err)
			http.Error(w, "Slackへの送信に失敗しました", http.StatusInternalServerError)
			return
		}

		// 送信したメッセージに対してスタンプ（リアクション）を自動で追加
		// ワークスペース独自のカスタムスタンプ（例: "tokyo-office"）なども指定可能
		reactionsToTarget := []string{"出勤", "退勤"}

		for _, reactionName := range reactionsToTarget {
			// slack.NewRefToMessage を使って「どのメッセージに押すか」を指定する
			msgRef := slack.NewRefToMessage(channelID, respTimestamp)

			err := api.AddReaction(reactionName, msgRef)
			if err != nil {
				// 1つスタンプが失敗してもメッセージ自体は飛んでいるので、ログだけ残して続行する
				log.Printf("スタンプ :%s: の追加に失敗: %v", reactionName, err)
			}
		}

		// 成功レスポンス
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"メッセージとスタンプを送信しました"}`))
	}
}
