package telegram

import (
	"time"
)

type Option func(o *options)
type options struct {
	apiUrlWithToken string
	interval        time.Duration
	timeout         time.Duration
	recvChanCap     uint
	sendChanCap     uint
}

func ApiUrlWithToken(url string) Option {
	return func(o *options) { o.apiUrlWithToken = url }
}

func Interval(i time.Duration) Option {
	return func(o *options) { o.interval = i }
}

func Timeout(t time.Duration) Option {
	return func(o *options) { o.timeout = t }
}

func RecvChanCap(cap uint) Option {
	return func(o *options) { o.recvChanCap = cap }
}

func SendChanCap(cap uint) Option {
	return func(o *options) { o.sendChanCap = cap }
}
