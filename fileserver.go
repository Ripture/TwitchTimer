package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ripture/TwitchTimer/lib"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
)

var GameList []forms.Games
var StreamerList []forms.Streamers
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
	//seed random with current time
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	//r.PathPrefix("/public/").Handler(http.FileServer(http.Dir("./public/")))
	r.PathPrefix("/public/").Handler(http.FileServer(http.Dir(".")))

	//websocket for requesting more streamers
	r.HandleFunc("/requestStreamer", requestStreamer)

	fmt.Printf("%v: Starting server on :1935\n", time.Now().Format(timeFormatString))
	http.ListenAndServe(":1935", r)
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

	// http.ServeFile(w, r, ".")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/js/main.js")
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

	fmt.Println(ip2)

	ip := serverIP{ip2}
	return ip
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	fname := path.Base(r.URL.Path)
	fmt.Println(r.URL.Path)
	fmt.Println(fname)
	http.ServeFile(w, r, "."+r.URL.Path)
}

func requestStreamer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in req")
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

		newStreamer := pickStreamer()

		fmt.Printf("%v: %v - Returning New Streamer: %v\n", time.Now().Format(timeFormatString), conn.RemoteAddr(), newStreamer)

		err = conn.WriteMessage(messageType, []byte(newStreamer))
		if err != nil {
			return

		}
	}
}

func pickStreamer() string {
	GetStreams()
	var ran = rand.Intn(len(StreamerList))
	return StreamerList[ran-1].Name
}

func GetStreams() forms.StreamS {
	var SomeStreams forms.StreamS
	var Streams forms.StreamS
	var offset int

	StreamerList = make([]forms.Streamers, 0)
	GameList = make([]forms.Games, 0)

	baseURL := "https://api.twitch.tv/kraken"

	resp, err := http.Get(baseURL + "/streams?game=Hearthstone%3A+Heroes+of+Warcraft")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &SomeStreams)
	if err != nil {
		panic(err)
	}

	//numt := int(float64(SomeStreams.Total)*.05) / 100 //number of times we need to get max limit (100) GETs
	//numm := int(float64(SomeStreams.Total)*.05) % 100 //number to pick up the remainder of limit for GET (will be limit=numm)
	//var skip = SomeStreams.Total - int(float64(SomeStreams.Total)*.05)
	numt := SomeStreams.Total / 100
	numm := SomeStreams.Total % 100
	fmt.Printf("%v: Found %v streamers\n", time.Now().Format("15:04:05AM"), SomeStreams.Total)

	for i := 0; i < numt; i++ {
		offset = 100 * i
		resp, err := http.Get(baseURL + "/streams?limit=100&game=Hearthstone%3A+Heroes+of+Warcraft&offset=" + strconv.Itoa(offset))
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &SomeStreams)
		if err != nil {
			panic(err)
		}

		Streams = BuildStreamList(Streams, SomeStreams)
	}

	offset = numt * 100

	if numm > 0 { //if there are any remaining streams to get
		resp, err := http.Get(baseURL + "/streams?limit=" + strconv.Itoa(numm) + "&game=Hearthstone%3A+Heroes+of+Warcraft&offset=" + strconv.Itoa(offset))
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &SomeStreams)
		if err != nil {
			panic(err)
		}

		Streams = BuildStreamList(Streams, SomeStreams)
	}

	return Streams
}

func BuildStreamList(Streams forms.StreamS, SomeStreams forms.StreamS) forms.StreamS {

	var tempGame forms.Games
	var tempStream forms.Streamers
	found := false
	_ = found

	for i := 0; i < len(SomeStreams.Streams); i++ {
		Streams.Streams = append(Streams.Streams, SomeStreams.Streams[i])
		tempGame.Name = SomeStreams.Streams[i].Game
		tempGame.Viewers = SomeStreams.Streams[i].Viewers
		tempStream.Viewers = tempGame.Viewers
		tempStream.Name = SomeStreams.Streams[i].Channel.Name
		tempStream.Game = SomeStreams.Streams[i].Game
		//if tempStream.Viewers > 0 {
		StreamerList = append(StreamerList, tempStream)
		//}

		for j := 0; j < len(GameList); j++ {
			if GameList[j].Name == tempGame.Name {
				GameList[j].Viewers += tempGame.Viewers
				found = true
			}

		}

		if !found {
			GameList = append(GameList, tempGame)
		}
	}

	return Streams
}
