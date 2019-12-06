package telegram

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maaxlee/kodi_youtube_bot/logger"
)

var (
	log = logger.GetLogger(os.Stdout, "Telegram: ", 0)
	ack bool
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
func RunTelBot(outCh chan string, ackCh chan bool, errorChan chan error) {

	log.Printf("Starting bot")
	token, ok := os.LookupEnv("TG_TOKEN")
	if !ok {
		errorChan <- fmt.Errorf("Token not found")
	}
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
	log.Debugp("Waiting for messages")
	for u := range updates {
		log.Debugp(fmt.Sprintf("Got message from channel: %s", u.Message.Text))
		id, err := getIdFromUrl(u.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Wrong or mailformed youtube URL, please check format")
			bot.Send(msg)
			continue
		}
		log.Debugp(fmt.Sprintf("Sending to channel id %s", id))
		outCh <- id
		ack = <-ackCh
		if !ack {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Something went wrong while trying to start playing, for details see logs")
			bot.Send(msg)
		}
	}
}
