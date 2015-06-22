package forms

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
