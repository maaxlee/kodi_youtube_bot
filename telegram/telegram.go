package telegram

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	logger = log.New(os.Stdout, "Telegram: ", 0)
	token  = ""
)

func getIdFromUrl(msg string) (string, error) {
	u, err := url.Parse(msg)
	if err != nil {
		return "", err
	}
	switch u.Host {
	case "www.youtube.com":
		q := u.Query()
		return q["v"][0], nil
	case "youtu.be":
		return strings.Trim(u.Path, "/"), nil
	default:
		return "", fmt.Errorf("Wrong or mailformed URL")
	}

}
func RunTelBot(ch chan string, errorChan chan error) {

	log.Print("Starting bot")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		errorChan <- err
	}
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60
	updates, err := bot.GetUpdatesChan(update)
	if err != nil {
		errorChan <- err
	}
	log.Print("Waiting for messages")
	for u := range updates {
		id, err := getIdFromUrl(u.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Wrong or mailformed youtube URL, please check format")
			bot.Send(msg)
			continue
		}
		log.Print("Sending to channel")
		ch <- id
	}
}
