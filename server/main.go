package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var xtokenFlag = flag.String("xtoken", "", "client xtoken from livecast")
var addr = flag.String("addr", "stagecast.se", "http service address")
var id = flag.String("id", "hypeisland", "project id")

type GameState struct {
	Type     string `json:"type"`
	MomentID string `json:"momentId"`
	NumUsers int    `json:"numUsers"`
	//	Users      []User               `json:"users"`
	Round      int                  `json:"round"`
	UserStates map[string]UserState `json:"userStates"`
}

type UserState map[string]string

type UserView struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
	View   string `json:"view"`
	Info   string `json:"info"`
}

var users = map[string]UserView{}
var gs GameState
var gsMutex sync.Mutex

func startTournament(momentID string) {
	gs = GameState{
		Type:     "server_info",
		MomentID: momentID,
		NumUsers: 10,
		//		Users:    []User{},
		Round: 0,
	}

}

func updateMomentGlobalState(momentID string) {
	log.Println(xtoken)

	u := url.URL{Scheme: "https", Host: "stagecast.se", Path: fmt.Sprintf("/api/moments/%s/state", momentID)}
	u = url.URL{Scheme: "http", Host: "localhost:8080", Path: fmt.Sprintf("/api/moments/state")}
	//https://stagecast.se/api/moments/{mid}/state©
	//https://stagecast.se/api/moments/{mid}/state
	userStates := map[string]UserState{}
	err := getJson(u.String(), &userStates)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Printf("%v", userStates)

	gsMutex.Lock()

	defer gsMutex.Unlock()

}

var xtoken string

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Token", xtoken)
	//req.Header.Set("Content-Type", "application/json")

	logrus.Debug("Send GET request: ", url)
	r, err := myClient.Do(req)
	logrus.Debug(r.StatusCode, r.Header, r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	flag.Parse()
	log.SetFlags(0)

	for _, userID := range []string{
		"1",
		"2",
		"3",
		"4",
	} {
		users[userID] = UserView{
			Type:   "client_info",
			UserID: userID,
			View:   "match",
			Info:   "This is a match! Don't lose.",
		}
	}

	xtoken = *xtokenFlag
	momentID := "60FCD373-2D61-47AB-B7DA-D2561DDFA66D"
	//updateMomentGlobalState(momentID)
	//log.Println("Exit")
	//os.Exit(0)

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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	currentView := "match"

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			log.Println("info")
			err := c.WriteJSON(gs)
			if err != nil {
				log.Println("write:", err)
				return
			}

			//Update state
			if currentView == "match" {
				currentView = "result"
			} else {
				currentView = "match"
			}

			for userID, view := range users {
				view.View = currentView
				log.Println(userID, view)

				err := c.WriteJSON(view)
				if err != nil {
					log.Println("write:", err)
					return
				}
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
