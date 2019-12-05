package main

import (
	"fmt"
)

type Request struct {
	id      int
	jsonrpc string
	method  string
	params  interface{}
}

type ClearParams struct {
	playlistid int
}

type AddParams struct {
	playlistid int
	item       struct {
		file string
	}
}

type OpenParams struct {
	item struct {
		playlistid int
		position   int
	}
}

var (
	clearRequest = Request{
		id:      755,
		jsonrpc: "2.0",
		method:  "Playlist.Clear",
		params:  ClearParams{playlistid: 1},
	}
	openRequest = Request{
		id:      756,
		jsonrpc: "2.0",
		method:  "Playlist.Open",
		params: OpenParams{item: struct {
			playlistid int
			position   int
		}{playlistid: 1, position: 0},
		},
	}
)

func main() {
	// {"id":752,"jsonrpc":"2.0","method":"Playlist.Clear","params":{"playlistid":1}}
	fmt.Println("Just started")
	fmt.Println(clearRequest)
	fmt.Println(openRequest)
}
