package profiles

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"
	"strings"

	"go.uber.org/zap"
)

// check how to make this across packages or in a struct. Copied from main right now
var ElgatoProgramPath = filepath.Join(os.Getenv("ProgramFiles"), "Elgato", "StreamDeck", "PageIcons")

func (p *ProfileList) getImage(firstState map[string]interface{}, pagePath string, UUID string) string {
	//	var button string
	image, hasProfileImg := firstState["Image"].(string)
	icon, hasElgatoImg := elgatoIcons[UUID]
	fmt.Printf(UUID)

	switch {
	case hasProfileImg:
		return p.fetchImg(image, pagePath)
	case hasElgatoImg:
		if icon.filename == "" {
			return ""
		}
		path := filepath.Join(ElgatoProgramPath, icon.path)
		return p.fetchImg(icon.filename, path)
	default:
		return ""
	}
}

func (p *ProfileList) fetchImg(imageUrl string, pagePath string) string {
	log := logger.GetDefaultLogger()
	path := filepath.Join(pagePath, imageUrl)
	img, err := os.ReadFile(path)
	if err != nil {
		log.Error("Error while reading image file",
			zap.String("path", path),
			zap.Error(err))
		return ""
	}
	var mime string
	switch strings.ToLower(filepath.Ext(imageUrl)) {
	case ".svg":
		mime = "image/svg+xml"
	default:
		mime = http.DetectContentType(img)
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(img)
}

var elgatoIcons = map[string]iconEntry{
	"com.elgato.streamdeck.multiactions.routine": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.page.goto": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.website": {
		path:     "com.elgato.streamdeck.pageicons.connectivity.sdIcons/Images/",
		filename: "10IconGlobe.svg",
	},

	"com.elgato.streamdeck.multiactions.routine2": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.page.indicator": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.hotkeyswitch": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.multiactions.random": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.soundboard.playaudio": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.hotkey": {
		path:     "com.elgato.streamdeck.pageicons.operations.sdIcons/Images/",
		filename: "17IconHotkey.svg",
	},

	"com.elgato.streamdeck.keys.logic": {
		path:     "com.elgato.streamdeck.pageicons.user.sdIcons/Images/",
		filename: "3IconTouch.svg",
	},

	"com.elgato.streamdeck.soundboard.stopaudioplay": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.open": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.profile.openchild": {
		path:     "com.elgato.streamdeck.pageicons.file.sdIcons/Images/",
		filename: "3IconFolder.svg",
	},

	"com.elgato.streamdeck.system.timer": {
		path:     "com.elgato.streamdeck.pageicons.tools.sdIcons/Images/",
		filename: "17IconTimer.svg",
	},

	"com.elgato.streamdeck.system.openapp": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.profile.rotate": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.keybrightness": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.close": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.page.previous": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "13IconChevronLeft.svg",
	},

	"com.elgato.streamdeck.system.sleep": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.text": {
		path:     "com.elgato.streamdeck.pageicons.textformatting.sdIcons/Images/",
		filename: "14IconText.svg",
	},

	"com.elgato.streamdeck.page.next": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "14IconChevronRight.svg",
	},

	"com.elgato.streamdeck.system.vsdtoggle": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},

	"com.elgato.streamdeck.system.multimedia": {
		path:     "com.elgato.streamdeck.pageicons.arrow.sdIcons/Images/",
		filename: "",
	},
}
