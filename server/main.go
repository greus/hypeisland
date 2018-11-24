package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "stagecast.se", "http service address")
var id = flag.String("id", "hypeisland", "project id")

type GameState struct {
	Type     string `json:"type"`
	MomentID string `json:"momentId"`
	NumUsers int    `json:"numUsers"`
	Users    []User `json:"users"`
	Round    int    `json:"round"`
}

type User struct {
}

func startTournament(momentID string) {
	gs = GameState{
		Type:     "server_info",
		MomentID: momentID,
		NumUsers: 10,
		Users:    []User{},
		Round:    0,
	}
}

var gs GameState

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: fmt.Sprintf("/api/events/%s/ws", *id)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	//sendChan := make(chan interface{})
	done := make(chan struct{})

	momentID := "60FCD373-2D61-47AB-B7DA-D2561DDFA66D"
	startTournament(momentID)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			log.Println("info")
			err := c.WriteJSON(gs)
			//			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
