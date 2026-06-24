package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// チャンネルの詳細情報を取得する
func HandleGetUserInfo(api *slack.Client, memberID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		users, err := api.GetUserInfo(memberID)
		if err != nil {
			log.Printf("users.info エラー: %v", err)
			http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// レンスポンス用のJSONを組み立てる
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"users":  users,
		})
	}
}
