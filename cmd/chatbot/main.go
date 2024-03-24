package main

import (
	"chatbot/clients/telegram"
	"chatbot/cmd/chatbot/config"
	"chatbot/db/mongo"
	tgprocess "chatbot/processes/telegram"
	"context"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	chatBot := &cli.App{
		Name:  "chatbot",
		Usage: "chat bot for telegram",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.yml",
				Usage: "config file path",
			},
		},
		Action: func(c *cli.Context) error {
			configFilePath := c.String("config")
			configBytes, err := os.ReadFile(configFilePath)
			if err != nil {
				log.Fatalln("read config failed:", err)
			}
			var chatBotConfig config.ChatBotConfig
			err = yaml.Unmarshal(configBytes, &chatBotConfig)
			if err != nil {
				log.Fatalln("parse config failed:", err)
			}

			apiUrlWithToken, err := url.JoinPath(chatBotConfig.TgApi.Url, chatBotConfig.TgApi.Token)
			if err != nil {
				log.Fatalln("url invalid:", err)
				return err
			}

			ctx := context.Background()

			mongo.InitMongoDB(chatBotConfig.Mongo.Uri)
			mongo.EnsureIndex(ctx)

			tgBot := telegram.New(
				telegram.ApiUrlWithToken(apiUrlWithToken),
				telegram.Interval(time.Duration(chatBotConfig.BotInfo.Interval)*time.Second),
				telegram.Timeout(time.Duration(chatBotConfig.BotInfo.HttpTimeout)*time.Second),
				telegram.RecvChanCap(chatBotConfig.BotInfo.RecvChanCap),
				telegram.SendChanCap(chatBotConfig.BotInfo.SendChanCap),
			)
			if err = tgBot.Init(ctx); err != nil {
				log.Fatalln("telegram client init failed:", err)
				return err
			}
			tgBot.Recv(context.Background(), tgprocess.DefaultCallback)

			return nil
		},
	}

	chatBot.Run(os.Args)
}
