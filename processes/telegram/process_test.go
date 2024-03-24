package telegram

import (
	"chatbot/clients/telegram"
	"chatbot/processes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_reviewHandler(t *testing.T) {
	type args struct {
		ctx context.Context
		req interface{}
	}
	nrStr, qrUrl, hlStr := "hello", "qrurl", "hotline"

	nr, _ := json.Marshal(telegram.Review{ContentType: telegram.TEXT, Text: &nrStr})
	qr, _ := json.Marshal(telegram.Review{ContentType: telegram.IMAGE, ImgUrl: &qrUrl})
	hotline, _ := json.Marshal(telegram.Review{ContentType: telegram.TEXT, Text: &hlStr})

	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"normal", args{ctx, nr}, nil, false},
		{"qrcode", args{ctx, qr}, "qrcode", false},
		{"hotline", args{ctx, hotline}, "hotline", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := processes.Chain(JsonMiddleware, ConvertRecordMiddleware)(reviewHandler)(tt.args.ctx, tt.args.req)
			if tt.name == "normal" {
				assert.Equal(t, resp, tt.want)
			} else {
				assert.Contains(t, string(resp.([]byte)), tt.want)
			}
		})
	}
}
