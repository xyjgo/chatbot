package telegram

type TGContentType int

const (
	TEXT TGContentType = iota
	IMAGE
	VOICE
	VIDEO
)

type Review struct {
	ReviewId    string        `json:"reviewid"`
	Sender      TGSender      `json:"sender"`
	Chat        TGChat        `json:"chat"`
	ContentType TGContentType `json:"contenttype"`
	Text        *string       `json:"text,omitempty"`
	ImgUrl      *string       `json:"imgurl,omitempty"`
	MediaUrl    *string       `json:"mediaurl,omitempty"`
	CreateTs    int64         `json:"createts"`
}

type TGSender struct {
	Uid       uint   `json:"uid"`
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatarurl"`
}

type TGChat struct {
	Cid      uint   `json:"cid"`
	ChatName string `json:"chatname"`
	ImgUrl   string `json:"imgurl"`
}
