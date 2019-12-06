package main

import (
	"github.com/maaxlee/kodi_youtube_bot/kodi"
	"github.com/maaxlee/kodi_youtube_bot/telegram"
)

func main() {

	ch := make(chan string, 5)
	errorChan := make(chan error, 2)
	go telegram.RunTelBot(ch, errorChan)
	go kodi.PlayYoutubeVideo(ch)
	for er := range errorChan {
		panic(er)
	}

}
