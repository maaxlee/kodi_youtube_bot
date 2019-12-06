package kodi

import (
	"math/rand"

	"github.com/gorilla/websocket"
)

type request struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type clearParams struct {
	Playlistid int `json:"playlistid"`
}
type addParams struct {
	Playlistid int     `json:"playlistid"`
	Item       addFile `json:"item"`
}

type addFile struct {
	File string `json:"file"`
}

type openParams struct {
	Item openItem `json:"item"`
}

type openItem struct {
	Playlistid int `json:"playlistid"`
	Position   int `json:"position"`
}

var (
	clearRequest = request{
		Id:      getId(),
		Jsonrpc: "2.0",
		Method:  "Playlist.Clear",
		Params:  clearParams{Playlistid: 1},
	}
	openRequest = request{
		Id:      getId(),
		Jsonrpc: "2.0",
		Method:  "Player.Open",
		Params: openParams{
			Item: openItem{
				Playlistid: 1,
				Position:   0,
			},
		},
	}
	url         = "ws://localhost:9090/jsonrpc"
	youtubePath = "plugin://plugin.video.youtube/play/?video_id="
)

func getId() int {
	return rand.Intn(1000)
}

func getAddRequest(youtubeId string) *request {
	return &request{
		Id:      getId(),
		Jsonrpc: "2.0",
		Method:  "Playlist.Add",
		Params: addParams{
			Playlistid: 1,
			Item:       addFile{File: youtubePath + youtubeId},
		},
	}
}

func playYoutubeVideo(videoId string) error {

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	err = ws.WriteJSON(clearRequest)
	if err != nil {
		return err
	}
	err = ws.WriteJSON(getAddRequest(videoId))
	if err != nil {
		return err
	}
	err = ws.WriteJSON(openRequest)
	if err != nil {
		return err
	}

	return nil
}

// Plays video on Kodi with given youtube video Id
func PlayYoutubeVideo(ch chan string) {
	for videoId := range ch {
		playYoutubeVideo(videoId)
	}
}
