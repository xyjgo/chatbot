package telegram

import (
	"chatbot"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	simserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer simserver.Close()
	ctx := context.Background()

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"init and stop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgBot := New(
				ApiUrlWithToken(simserver.URL),
				Interval(time.Duration(1)*time.Second),
				Timeout(time.Duration(2)*time.Second),
				RecvChanCap(5),
				SendChanCap(5),
			)

			err := tgBot.Init(ctx)
			assert.Nil(t, err)

			err = tgBot.Send(ctx, []byte{0})
			assert.Nil(t, err)

			recvCnt := 0
			recvClosed := false
			go func() {
				tgBot.Recv(ctx, func(client chatbot.Client, rawData []byte) {
					recvCnt++
				})
				recvClosed = true
			}()

			time.Sleep(2 * time.Second)
			assert.Greater(t, recvCnt, 0)

			tgBot.Stop(ctx)
			time.Sleep(2 * time.Second)

			err = tgBot.Send(ctx, []byte{0})
			assert.NotNil(t, err)
			assert.Equal(t, recvClosed, true)

		})
	}
}
