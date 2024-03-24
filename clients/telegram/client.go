package telegram

import (
	"bytes"
	"chatbot"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	opts     *options
	client   *http.Client
	ctx      context.Context
	cancel   context.CancelFunc
	recvChan chan []byte
	sendChan chan []byte
}

// New make a client for telegram api
func New(opts ...Option) *Client {
	o := options{
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(&o)
	}

	client := &http.Client{
		Timeout: o.timeout,
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{&o, client, ctx, cancel, make(chan []byte, o.recvChanCap), make(chan []byte, o.sendChanCap)}
}

// Init start sendChan consumer routine, do prepare job like login and etc
func (c *Client) Init(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case r := <-c.sendChan:
				c.SendSync(ctx, r)
			}
		}
	}()
	return nil
}

// Stop  cancel send and recv routine and action
func (c *Client) Stop(ctx context.Context) error {
	c.cancel()
	return nil
}

// Send put data into sendChan
func (c *Client) Send(ctx context.Context, rawData []byte) error {
	select {
	case <-c.ctx.Done():
		return errors.New("client closed")
	default:
		c.sendChan <- rawData
	}

	return nil
}

// SendSync send response to third part api
func (c *Client) SendSync(ctx context.Context, rawData []byte) ([]byte, error) {
	resp, err := c.client.Post(c.opts.apiUrlWithToken, "application/json", bytes.NewReader(rawData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

// Recv continuous receiving message until stopped
func (c *Client) Recv(ctx context.Context, callback chatbot.Callback) error {
	// start recvChan consumer routine, get review from recvChan and callback
	select {
	case <-c.ctx.Done():
		return errors.New("client closed")
	default:
		go func() {
			for {
				select {
				case <-c.ctx.Done():
					return
				case r := <-c.recvChan:
					if callback != nil {
						callback(c, r)
					}
				}
			}
		}()
	}

	// continuous receiving review, put into recvChan

	for {
		select {
		case <-c.ctx.Done():
			return errors.New("client closed")
		case <-time.After(c.opts.interval):
		}

		resp, err := c.client.Get(c.opts.apiUrlWithToken)
		if err != nil {
			log.Println("recv err:", err)
			continue
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		c.recvChan <- body
	}
	return nil
}

// NewRecordFromTgReview convert Review to ReviewRecord format
func NewRecordFromTgReview(review *Review) *chatbot.ReviewRecord {
	record := chatbot.ReviewRecord{
		Vendor:          chatbot.TELEGRAM,
		ReviewId:        review.ReviewId,
		SenderUid:       review.Sender.Uid,
		SenderNickname:  review.Sender.Nickname,
		SenderAvatarUrl: review.Sender.AvatarUrl,
		ChatId:          review.Chat.Cid,
		ChatName:        review.Chat.ChatName,
		ChatImgUrl:      review.Chat.ImgUrl,
		ContentType:     int(review.ContentType),
		CreateTs:        review.CreateTs,
	}
	if review.Text != nil {
		record.Text = *review.Text
	}
	if review.ImgUrl != nil {
		record.ImgUrl = *review.ImgUrl
	}
	if review.MediaUrl != nil {
		record.MediaUrl = *review.MediaUrl
	}
	return &record
}

func NewReviewFromRecord(record *chatbot.ReviewRecord) *Review {
	review := Review{
		ReviewId: record.ReviewId,
		Sender: TGSender{
			record.SenderUid, record.SenderNickname, record.SenderAvatarUrl,
		},
		Chat:        TGChat{record.ChatId, record.ChatName, record.ChatImgUrl},
		ContentType: TGContentType(record.ContentType),
		Text:        &record.Text,
		ImgUrl:      &record.ImgUrl,
		MediaUrl:    &record.MediaUrl,
		CreateTs:    record.CreateTs,
	}
	return &review
}
