package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// チャンネル内のメンバーID一覧を取得する
func HandleGetMembers(api *slack.Client, channelID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Slack API: conversations.members
		// 1回で取得する上限を100名に設定（必要に応じて増やせる）
		params := &slack.GetUsersInConversationParameters{
			ChannelID: channelID,
			Limit:     100,
		}

		memberIDs, _, err := api.GetUsersInConversation(params)
		if err != nil {
			log.Printf("conversations.members エラー: %v", err)
			http.Error(w, "メンバー一覧の取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// レンスポンス用のJSONを組み立てる
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"members": memberIDs,
		})
	}
}
