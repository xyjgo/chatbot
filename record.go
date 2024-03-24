package chatbot

const (
	TELEGRAM = "telegram"
	WECHAT   = "wechat"
	QQ       = "qq"
)

type ReviewRecord struct {
	Vendor    string
	ReviewId  string
	ReplyToId string

	SenderUid       uint
	SenderNickname  string
	SenderAvatarUrl string

	ChatId     uint
	ChatName   string
	ChatImgUrl string

	ContentType int
	Text        string
	ImgUrl      string
	MediaUrl    string
	CreateTs    int64
}
