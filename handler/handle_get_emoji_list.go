package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// ワークスペース内のカスタムスタンプ一覧を取得する
func HandleGetEmojiList(api *slack.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Slack API: emoji.list
		emojiMap, err := api.GetEmoji()
		if err != nil {
			log.Printf("emoji.list エラー: %v", err)
			http.Error(w, "絵文字一覧の取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// マップの「キー（スタンプ名）」だけを格納する配列（スライス）を用意
		emojiNames := make([]string, 0, len(emojiMap))
		for name := range emojiMap {
			emojiNames = append(emojiNames, name)
		}

		// emojiMap の中身は map[string]string (例: "tokyo-office": "https://...") になっている
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"emojis": emojiNames,
		})
	}
}
