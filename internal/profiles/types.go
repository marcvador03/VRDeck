package profiles

import "VRDeck/internal/ws"

type Buttons struct {
	Row       int    `json:"row"`
	Col       int    `json:"col"`
	Title     string `json:"label"`
	ActionID  string `json:"-"`
	Icon      string `json:"icon,omitempty"`
	IconType  string `json:"icontype,omitempty"`
	UUID      string `json:"-"`
	ChildPage *Pages `json:"-"`
}

type Pages struct {
	UUID     string    `json:"-"`
	Name     string    `json:"-"`
	Buttons  []Buttons `json:"buttons"`
	pagePath string    `json:"-"`
	//parentPage *Pages    `json:"-"`
}

type Profile struct {
	UUID     string
	Name     string
	pagesNum int
	Pages    []*Pages
	Default  *Pages
	Current  *Pages
}

type ProfileList struct {
	path     string
	ws       *ws.MSFSWebSocket
	Profiles []Profile
	Current  string
}

type CurrentPage struct {
	Pages struct {
		Current string `json:"Current"`
	} `json:"Pages"`
}

type iconEntry struct {
	path     string
	filename string
}
