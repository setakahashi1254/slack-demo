package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// チャンネルの詳細情報を取得する
func HandleGetChannelInfo(api *slack.Client, channelID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		params := &slack.GetConversationInfoInput{
			ChannelID:         channelID,
			IncludeLocale:     true,
			IncludeNumMembers: true,
		}

		conversations, err := api.GetConversationInfo(params)
		if err != nil {
			log.Printf("conversations.info エラー: %v", err)
			http.Error(w, "チャンネル詳細情報の取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// レンスポンス用のJSONを組み立てる
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "success",
			"conversations": conversations,
		})
	}
}
