package telegram

import (
	"chatbot"
	"chatbot/clients/telegram"
	"chatbot/db/mongo"
	"chatbot/processes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

func JsonMiddleware(handler processes.Handler) processes.Handler {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		var review telegram.Review
		json.Unmarshal(req.([]byte), &review)
		resp, err = handler(ctx, &review)
		if resp != nil {
			resp, err = json.Marshal(resp)
		}
		return
	}
}

func LogMiddleware(handler processes.Handler) processes.Handler {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		log.Printf("recv %v\n", string(req.([]byte)))
		resp, err = handler(ctx, req)
		if resp != nil {
			log.Printf("send %v\n", string(resp.([]byte)))
		}
		return
	}
}

func ConvertRecordMiddleware(handler processes.Handler) processes.Handler {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		review := req.(*telegram.Review)
		record := telegram.NewRecordFromTgReview(review)
		resp, err = handler(ctx, record)
		if resp != nil {
			resp = telegram.NewReviewFromRecord(resp.(*chatbot.ReviewRecord))
		}
		return
	}
}

func recordDbMiddleware(handler processes.Handler) processes.Handler {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		record := req.(*chatbot.ReviewRecord)
		mongo.SaveReviewRecord(ctx, record)
		resp, err = handler(ctx, record)
		if resp != nil {
			mongo.SaveReviewRecord(ctx, resp.(*chatbot.ReviewRecord))
		}
		return
	}
}

func reviewHandler(ctx context.Context, req interface{}) (interface{}, error) {
	record := req.(*chatbot.ReviewRecord)

	handles := []processes.Handle{processes.HandleQrcode, processes.HandleHotLine}

	var content string
	for _, handle := range handles {
		if resp, _ := handle(record); resp != nil {
			content = string(resp)
			break
		}
	}

	if len(content) == 0 {
		return nil, nil
	}

	reply := chatbot.ReviewRecord{
		ReviewId:        uuid.NewString(),
		ReplyToId:       record.ReviewId,
		SenderUid:       0,
		SenderNickname:  "chatbot",
		SenderAvatarUrl: "avatarurl",
		ChatId:          record.ChatId,
		ChatName:        record.ChatName,
		ChatImgUrl:      record.ChatImgUrl,
		ContentType:     0,
		Text:            content,
		CreateTs:        time.Now().Unix(),
	}
	return &reply, nil
}

func DefaultCallback(client chatbot.Client, rawData []byte) {
	resp, _ := processes.Chain(LogMiddleware, JsonMiddleware, ConvertRecordMiddleware, recordDbMiddleware)(reviewHandler)(context.Background(), rawData)
	if resp != nil {
		client.Send(context.Background(), resp.([]byte))
	}
}
