package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/moments/state", States)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func States(w http.ResponseWriter, r *http.Request) {

	names := []string{
		"angelina",
		"magnus",
		"sebastian",
		"johan",
	}
	users := []User{}
	for i, name := range names {

		user := User{
			State: map[string]string{
				"score": strconv.Itoa(i),
				"name":  name,
			},
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

type User struct {
	State map[string]string `json:"state"`
}
