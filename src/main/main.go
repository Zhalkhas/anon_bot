package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	token = "<Here goes token of telegram bot>"
	id    = -0000000000 // Here goes ID of chat
	url   = "<Here goes URL of server>"
)

var (
	bot *tgbotapi.BotAPI
)

func initTelegram() {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Bot init error %s", err.Error())
	}
	info, err := bot.GetWebhookInfo()
	log.Println("Webhook info:", info)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url + token))
	if err != nil {
		log.Printf("Webhook setting error: %s", err.Error())
	} else {
		log.Println("Webhook set successfully")
	}
}

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	if update.Message != nil {
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "Just write your message and I will resend it anonymously."
			case "start":
				msg.Text = "Just write your message and I will resend it anonymously."
			default:
				msg.Text = "Unknown command, try /help"
			}
			bot.Send(msg)
		} else if update.Message.Chat.ID != id {

			msg := tgbotapi.NewMessage(id, update.Message.Text)

			if update.Message.ForwardFrom != nil && update.Message.ForwardFromChat != nil {
				bot.Send(tgbotapi.NewForward(id, update.Message.ForwardFromChat.ID, update.Message.ForwardFromMessageID))
			}

			if update.Message.Photo != nil {
				photo := *update.Message.Photo
				fileid := photo[len(photo)-1].FileID
				if update.Message.Caption != "" {
					msg.Text = update.Message.Caption
				}
				bot.Send(tgbotapi.NewPhotoShare(id, fileid))
			}

			if update.Message.Sticker != nil {
				bot.Send(tgbotapi.NewStickerShare(id, update.Message.Sticker.FileID))
				if update.Message.Caption != "" {
					msg.Text = update.Message.Caption
				}
			}

			if update.Message.Document != nil {
				bot.Send(tgbotapi.NewDocumentShare(id, update.Message.Document.FileID))
				if update.Message.Caption != "" {
					msg.Text = update.Message.Caption
				}
			}

			if update.Message.Video != nil {
				bot.Send(tgbotapi.NewVideoShare(id, update.Message.Video.FileID))
				if update.Message.Caption != "" {
					msg.Text = update.Message.Caption
				}
			}
			if update.Message.VideoNote != nil {
				bot.Send(tgbotapi.NewVideoNoteShare(id, update.Message.VideoNote.Length, update.Message.VideoNote.FileID))
				if update.Message.Caption != "" {
					msg.Text = update.Message.Caption
				}
			}

			bot.Send(msg)
		}
	} else {
		log.Printf("Empty message recieved %+v", update)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Printf("Port set to :%s", port)
	// gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	log.Println("Router initialized")

	// telegram
	initTelegram()
	log.Printf("Telegram API init sequence complete")
	router.POST("/"+token, webhookHandler)
	log.Println("Listening for /<token> started")
	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
