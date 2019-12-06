package kodi

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/maaxlee/kodi_youtube_bot/logger"
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

type response struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

var (
	url         = fmt.Sprintf("ws://%s:9090/jsonrpc", getEnv("KODI_HOST", "localhost"))
	youtubePath = "plugin://plugin.video.youtube/play/?video_id="
	log         = logger.GetLogger(os.Stdout, "Kodi: ", 0)
)

func getEnv(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}

func getId() int {
	return rand.Intn(1000)
}

func getClearRequest() (*request, int) {
	id := getId()
	clearRequest := &request{
		Id:      id,
		Jsonrpc: "2.0",
		Method:  "Playlist.Clear",
		Params:  clearParams{Playlistid: 1},
	}
	return clearRequest, id

}

func getOpenRequest() (*request, int) {
	id := getId()
	openRequest := &request{
		Id:      id,
		Jsonrpc: "2.0",
		Method:  "Player.Open",
		Params: openParams{
			Item: openItem{
				Playlistid: 1,
				Position:   0,
			},
		},
	}
	return openRequest, id

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

func checkResponse(ws *websocket.Conn, req string, id int) error {

	for {
		select {
		case <-time.After(10 * time.Second):
			return fmt.Errorf("Could not wait until successfull response for request %s", req)
		default:
			r := &response{}
			err := ws.ReadJSON(r)
			if err != nil {
				return err
			}
			if r.Id != id {
				continue
			}
			if r.Result != "OK" {
				return fmt.Errorf("Non OK response after %s request: %v", req, r)
			}
			return nil
		}
	}

}
func playYoutubeVideo(videoId string) error {

	log.Debugp("Opening WS to kodi")
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	log.Debugp("Sending Clear request")
	clearRequest, id := getClearRequest()
	err = ws.WriteJSON(clearRequest)
	if err != nil {
		return err
	}
	if err := checkResponse(ws, "clear", id); err != nil {
		return err
	}

	log.Debugp("Sending add video request")
	err = ws.WriteJSON(getAddRequest(videoId))
	if err != nil {
		return err
	}

	log.Debugp("Sending open request")
	openRequest, id := getOpenRequest()
	err = ws.WriteJSON(openRequest)
	if err != nil {
		return err
	}

	if err := checkResponse(ws, "open", id); err != nil {
		return err
	}

	return nil
}

// Plays video on Kodi with given youtube video Id
func PlayYoutubeVideo(playCh chan string, ackCh chan bool) {
	for videoId := range playCh {
		log.Debugp("Got video Id to play")
		err := playYoutubeVideo(videoId)
		if err != nil {
			log.Debugp("Error on sending data to Kodi")
			log.Debugp(err)
			ackCh <- false
			continue
		}
		ackCh <- true
	}
}
