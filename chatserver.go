package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Use capitalized filed name to enable json marshalling
// of field, and then
// controll json dump use smaller case field name
type Message struct {
	Timestamp float64 `json:"timestamp"`
	User      string  `json:"user"`
	Text      string  `json:"text"`
}

// A slice of capacity 200 as internal storage to hold
// all messages need to be queried.
const capacity = 100

var messages = make([]Message, 0, capacity*2)

// Use boolean map serve as a set to hold all users
var users = make(map[string]bool)

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
	message := Message{}
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
	message.Timestamp = getCurtimeStamp()
	fmt.Println("user:" + message.User)
	fmt.Println("text:" + message.Text)

	// Store message and user
	if len(messages) == cap(messages) {
		z := make([]Message, len(messages), capacity*2)
		copy(z, messages)
		messages = z
	}
	if len(messages) == capacity {
		messages = messages[1:]
	}
	messages = append(messages, message)

	// save user simply without concerning about capacity
	users[message.User] = true

	result := make(map[string]bool)
	result["ok"] = true
	encoder := json.NewEncoder(c)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

// this logic is to controll output format,
// using second as granuality, with two more digits
func getCurtimeStamp() float64 {
	t := time.Now()
	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	return float64(tUnixMilli/10) / 100
}

// list all users, this is simple
func listUsers(c http.ResponseWriter, req *http.Request) {
	// var data []byte
	// userlist := UserList{}
	var userl = make([]string, 0, len(users))
	for k := range users {
		userl = append(userl, k)
		fmt.Println("find user:" + k)
	}
	result := make(map[string][]string)
	result["users"] = userl[:]
	encoder := json.NewEncoder(c)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

// list messages, this is also simple
func listMessages(c http.ResponseWriter, req *http.Request) {
	result := make(map[string][]Message)
	result["messages"] = messages[:]
	encoder := json.NewEncoder(c)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}
