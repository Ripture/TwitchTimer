package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	GetStreams()

	r := mux.NewRouter()

	r.HandleFunc("/", Twitch)

	r.HandleFunc("/next", Next)

	fmt.Println("Starting server on :1935")
	http.ListenAndServe(":1935", r)
}

func Twitch(w http.ResponseWriter, r *http.Request) {
	var ran = rand.Intn(len(StreamerList))
	fmt.Println("Fetching " + StreamerList[ran].Name)
	ExeTemplate(StreamerList[ran].Name, StreamerList[ran].Viewers, w, r)
}

func Next(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching new streamer list...")
	GetStreams()
	var ran = rand.Intn(len(StreamerList))
	fmt.Println("Fetching " + StreamerList[ran].Name)
	ExeTemplate(StreamerList[ran].Name, StreamerList[ran].Viewers, w, r)
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
	fmt.Println("Found " + strconv.Itoa(SomeStreams.Total) + " streamers.")

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
