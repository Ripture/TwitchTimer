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

type Streamer struct {
	Name    string
	Viewers int
}

// WHAT THE HELL IS THIS
//
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	GetStreams()
	r := mux.NewRouter()

	r.HandleFunc("/", Twitch)

	r.HandleFunc("/requestStreamer", requestStreamer)

	fmt.Printf("%v: Starting server on :1935\n", time.Now().Format("15:04:05AM"))
	http.ListenAndServe(":1935", r)
}

func requestStreamer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		_ = messageType
		_ = p
		if err != nil {
			return
		}

		msg := string(p[:])
		// fmt.Println(msg)
		if msg == "requestStreamer" {
			fmt.Printf("%v: %v - Requests New Streamer\n", time.Now().Format("15:04:05AM"), conn.RemoteAddr())

			newStreamer := pickStreamer()

			fmt.Printf("%v: %v - Returning New Streamer: %v\n", time.Now().Format("15:04:05AM"), conn.RemoteAddr(), newStreamer)

			err = conn.WriteMessage(messageType, []byte(newStreamer))
			if err != nil {
				return
			}
		}
	}
}
func print_binary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%c", s[n])
	}
	fmt.Printf("\n")
}

func Twitch(w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("%v: New Connection From %v\n", time.Now().Format("15:04:05AM"), conn.RemoteAddr())

	var ran = rand.Intn(len(StreamerList))
	ExeTemplate(StreamerList[ran].Name, StreamerList[ran].Viewers, w, r)
}

func pickStreamer() string {
	//fmt.Println("Fetching new streamer list...")
	GetStreams()
	var ran = rand.Intn(len(StreamerList))
	//fmt.Println("Fetching " + StreamerList[ran-1].Name)
	return StreamerList[ran-1].Name
}

func ExeTemplate(name string, viewers int, w http.ResponseWriter, r *http.Request) {
	stre := Streamer{name, viewers}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, stre); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
