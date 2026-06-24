package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slack-demo/handler"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// 環境変数からトークンを取得
	botToken := os.Getenv("SLACK_BOT_TOKEN") // xoxb-...
	appToken := os.Getenv("SLACK_APP_TOKEN") // xapp-...

	// 環境変数からチャンネルIDを取得
	channelID := os.Getenv("CHANNEL_ID") // C...

	// 環境変数からメンバーIDを取得
	memberID := os.Getenv("MEMBER_ID") // U...

	if botToken == "" || appToken == "" {
		log.Fatal("SLACK_BOT_TOKEN と SLACK_APP_TOKEN を設定してください")
	}

	// クライアントの初期化
	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(api)

	// ルーティング登録
	http.HandleFunc("/send-message", handler.HandlePostMessage(api, channelID))
	http.HandleFunc("/get-members", handler.HandleGetMembers(api, channelID))
	http.HandleFunc("/get-channel-info", handler.HandleGetChannelInfo(api, channelID))
	http.HandleFunc("/get-user-info", handler.HandleGetUserInfo(api, memberID))
	http.HandleFunc("/send-ephemeral", handler.HandlePostEphemeral(api, channelID))
	http.HandleFunc("/get-emojis", handler.HandleGetEmojiList(api))
	http.HandleFunc("/send-message-with-stamp", handler.HandlePostMessageWithStamp(api, channelID))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// イベント監視を裏で実行
	go handler.HandleSlackEvents(ctx, client)

	// Socket Modeの起動自体も裏で実行する
	go func() {
		fmt.Println("Socket Modeで常時監視を開始しました...")
		err := client.RunContext(ctx)
		if err != nil {
			log.Fatalf("Socket Modeの起動に失敗しました: %v", err)
		}
	}()

	// Dockerコンテナ内から外に公開するため、"0.0.0.0:8080" で待ち受ける
	fmt.Println("Goサーバーがポート8080で起動しました。リクエストを待っています...")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}
