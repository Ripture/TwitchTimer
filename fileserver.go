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
	Name string
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
	var streamer = StreamerList[rand.Intn(len(StreamerList))].Name
	ExeTemplate(streamer, w, r)
}

func Next(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting another random streamer")
	var streamer = StreamerList[rand.Intn(len(StreamerList))].Name
	fmt.Println("fetching " + streamer)
	ExeTemplate(streamer, w, r)
}

func ExeTemplate(name string, w http.ResponseWriter, r *http.Request) {
	stre := Streamer{name}

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

	baseURL := "https://api.twitch.tv/kraken"

	resp, err := http.Get(baseURL + "/streams?limit=100")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &SomeStreams)
	if err != nil {
		panic(err)
	}

	numt := int(float64(SomeStreams.Total)*.05) / 100 //number of times we need to get max limit (100) GETs
	numm := int(float64(SomeStreams.Total)*.05) % 100 //number to pick up the remainder of limit for GET (will be limit=numm)
	var skip = SomeStreams.Total - int(float64(SomeStreams.Total)*.05)

	for i := 0; i < numt; i++ {
		offset = 100 * i
		resp, err := http.Get(baseURL + "/streams?limit=100&offset=" + strconv.Itoa(offset+skip))
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
		resp, err := http.Get(baseURL + "/streams?limit=" + strconv.Itoa(numm) + "&offset=" + strconv.Itoa(offset+skip))
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
		StreamerList = append(StreamerList, tempStream)

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
