package profiles

import "streamdeckVR/internal/ws"

type Buttons struct {
	Row      int    `json:"row"`
	Col      int    `json:"col"`
	Title    string `json:"label"`
	ActionID string `json:"-"`
}

type Pages struct {
	UUID    string    `json:"-"`
	Name    string    `json:"-"`
	Buttons []Buttons `json:"buttons"`
}

type Profile struct {
	UUID     string
	Name     string
	pagesNum int
	Pages    []*Pages
}

type ProfileList struct {
	path     string
	ws       *ws.MSFSWebSocket
	Profiles []Profile
}

type CurrentPage struct {
	Pages struct {
		Current string `json:"Current"`
	} `json:"Pages"`
}
