package chatbot

import "context"

type Callback func(client Client, rawData []byte)

type Client interface {
	Init(ctx context.Context) error
	Stop(ctx context.Context) error
	Send(ctx context.Context, rawData []byte) error
	Recv(ctx context.Context, callback Callback) error
}
