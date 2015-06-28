package forms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PreviewS struct {
	Small    string
	Medium   string
	Large    string
	Template string
}

type ChannelLinksS struct {
	Self          string
	Follows       string
	Commercial    string
	StreamKey     string `json:"stream_key"`
	Chat          string
	Features      string
	Subscriptions string
	Editors       string
	Teams         string
	Videos        string
}

type ChannelAttrS struct {
	Mature               bool
	Status               string
	BroadcasterLang      string `json:"broadcaster_language"`
	DisplayName          string `json:"display_name"`
	Game                 string
	Delay                int
	Language             string
	ID                   int `json:"_id"`
	Name                 string
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	Logo                 string
	Banner               string
	VideoBanner          string `json:"video_banner"`
	Background           string
	ProfileBanner        string `json:"profile_banner"`
	ProfileBannerBGColor string `json:"profile_banner_background_color"`
	Partner              bool
	URL                  string
	Views                int
	Followers            int
	Links                ChannelLinksS `json:"_links"`
}

type ChannelS struct {
	Game      string
	Viewers   int
	CreatedAt string `json:"created_at"`
	ID        int    `json:"_id"`
	Channel   ChannelAttrS
	Preview   PreviewS
	Links     LinkS `json:"_links"`
}

type LinkS struct {
	Summary  string
	Followed string
	Next     string
	Featured string
	Self     string
}

type StreamS struct {
	Total   int `json:"_total"`
	Streams []ChannelS
	Links   LinkS `json:"_links"`
}

type Streamers struct {
	Name    string
	Game    string
	Viewers int
}

type Games struct {
	Name    string
	Viewers int
}

func GetStreams() StreamS {
	var SomeStreams StreamS
	var Streams StreamS
	var offset int

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

	numt := SomeStreams.Total / 100
	numm := SomeStreams.Total % 100

	for i := 0; i < SomeStreams.Total/100; i++ {
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

		Streams.Streams = append(Streams.Streams, SomeStreams.Streams...)
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

		Streams.Streams = append(Streams.Streams, SomeStreams.Streams...)
	}

	return Streams
}
