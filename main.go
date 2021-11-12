package main

import (
	"os"

	"github.com/maaxlee/kodi_youtube_bot/kodi"
	"github.com/maaxlee/kodi_youtube_bot/logger"
	"github.com/maaxlee/kodi_youtube_bot/telegram"
)

var log = logger.GetLogger(os.Stdout, "Main: ", 0)

func main() {

	log.Printf("Starting telegram youtube bot")
	tubeCh := make(chan string, 5)
	torCh := make(chan string, 5)
	ackCh := make(chan bool, 5)
	errorChan := make(chan error, 2)
	go telegram.RunTelBot(tubeCh, torCh, ackCh, errorChan)
	go kodi.HandleKodiInput(tubeCh, torCh, ackCh)
	for err := range errorChan {
		panic(err)
	}

}
