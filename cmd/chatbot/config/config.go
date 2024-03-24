package config

type ChatBotConfig struct {
	TgApi   TelegramApi `yaml:"tg_api"`
	BotInfo ChatBot     `yaml:"bot_info"`
	Mongo   MongoDb     `yaml:"mongo"`
}

type TelegramApi struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type ChatBot struct {
	Interval    uint `yaml:"interval"`
	HttpTimeout uint `yaml:"http_timeout"`
	RecvChanCap uint `yaml:"recv_chan_cap"`
	SendChanCap uint `yaml:"send_chan_cap"`
}

type MongoDb struct {
	Uri string `yaml:"uri"`
}
