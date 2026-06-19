package handler

import (
	"context"
	"fmt"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// HandleSlackEvents - Socket Mode経由で届くSlackのイベントを常時監視する
func HandleSlackEvents(ctx context.Context, client *socketmode.Client) {
	for evt := range client.Events {
		// アプリ終了時（Contextキャンセル時）に安全にループを抜けるお作法
		if ctx.Err() != nil {
			return
		}

		switch evt.Type {
		// SlackのEvent APIからデータ（スタンプなど）が届いた場合
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}

			// 必須：Slack側に「通信を受け取ったよ」と即座に応答（Ack）を返す
			client.Ack(*evt.Request)

			// 具体的なイベントの中身を判定
			switch eventsAPIEvent.InnerEvent.Type {
			case "reaction_added": // スタンプが押されたイベント
				ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
				if !ok {
					continue
				}
				fmt.Printf("【スタンプ追加】ユーザー: %s | スタンプ: :%s: | 対象メッセージTS: %s\n", ev.User, ev.Reaction, ev.Item.Timestamp)

			case "reaction_removed": // スタンプが消されたイベント
				ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionRemovedEvent)
				if !ok {
					continue
				}
				fmt.Printf("【スタンプ削除】ユーザー: %s | スタンプ: :%s: | 対象メッセージTS: %s\n", ev.User, ev.Reaction, ev.Item.Timestamp)

			case "member_joined_channel": // チャンネルにユーザーが参加したイベント
				ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MemberJoinedChannelEvent)
				if !ok {
					continue
				}
				fmt.Printf("【ユーザー参加】ユーザー: %s | チャンネル: %s\n", ev.User, ev.Channel)
			}

		}
	}
}
