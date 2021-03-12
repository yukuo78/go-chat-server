package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	timestamp float64
	user      string
	text      string
}

// A slice of capacity 100
var Messages []Message

// var UserSet := map[string]bool
var users []string

func main() {
	http.HandleFunc("/status",
		func(c http.ResponseWriter, req *http.Request) {
			c.Write([]byte("alive"))
		})

	http.HandleFunc("/message", postMessage)
	http.HandleFunc("/users", listUsers)
	http.HandleFunc("/messages", listMessages)
	err := http.ListenAndServe(":8081", nil)
	panic(err)
}

// read message content from httprequest as a json
// extract: user, message, timestamp
// save to localstorage, only recent 100 messages is needed
// so we can store the messages in a list
func postMessage(c http.ResponseWriter, req *http.Request) {
	// save message
	// save user
	message := make(map[string]string)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("reqd request body failed!")
	}
	fmt.Println(reqBody)
	// fmt.Println("user:" + message["user"])
	errb := json.Unmarshal(reqBody, &message)
	if errb != nil {
		log.Fatal("read body failed!")
	}
	fmt.Println("user:" + message["user"])
	fmt.Println("text:" + message["text"])
	// if len(Messages) == 100 {
	// 	Messages = Messages[1:]
	// }
	// append(messages, message)
	// UserSet[message.user] = true
	users = append(users, message["user"])

	result := make(map[string]bool)
	result["ok"] = true
	var data []byte
	data, _ = json.Marshal(result)
	fmt.Fprintf(c, "%s\n", data)
}

// list all users, this is simple
func listUsers(c http.ResponseWriter, req *http.Request) {
	var data []byte
	data, _ = json.Marshal(users)
	fmt.Fprintf(c, "%s", data)
}

// list messages, this is also simple
func listMessages(c http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(c, "[hello]")
}
