package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ripture/TwitchTimer/lib"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var timeFormatString = "15:04:05AM"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type serverIP struct {
	Ip string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1935"
	}

	//seed random with current time
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.PathPrefix("/public/").Handler(http.FileServer(http.Dir(".")))

	//websocket for requesting more streamers
	r.HandleFunc("/requestStreamer", requestStreamer)

	fmt.Printf("%v: Starting server on :1935\n", time.Now().Format(timeFormatString))
	http.ListenAndServe(":"+port, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ip := getServerIP()

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getServerIP() serverIP {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	ip1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	ip2 := (string(ip1))[:len(string(ip1))-1]

	ip := serverIP{ip2}
	return ip
}

func requestStreamer(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("%v: %v - Requests New Streamer\n", time.Now().Format(timeFormatString), conn.RemoteAddr())

		Streams := forms.GetStreams()
		newStreamer := Streams.Streams[rand.Intn(len(Streams.Streams)-1)].Channel.DisplayName

		fmt.Printf("%v: %v - Returning New Streamer: %v\n", time.Now().Format(timeFormatString), conn.RemoteAddr(), newStreamer)

		err = conn.WriteMessage(messageType, []byte(newStreamer))
		if err != nil {
			return

		}
	}
}
