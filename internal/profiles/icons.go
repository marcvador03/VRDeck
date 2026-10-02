package profiles

import (
	"VRDeck/internal/logger"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func (p *Profile) getImage(firstState map[string]interface{}, pagePath string, UUID string) (string, string) {
	//	var button string
	image, hasProfileImg := firstState["Image"].(string)
	icon, hasElgatoImg := elgatoIcons[UUID]
	switch {
	case hasProfileImg:
		return p.fetchImg(image, pagePath), ""
	case hasElgatoImg:
		if icon.filename == "" {
			return "", ""
		}
		return p.fetchImg(icon.filename, icon.path), "half"
	default:
		return "", ""
	}
}

func (p *Profile) fetchImg(imageUrl string, pagePath string) string {
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

// https://github.com/czottmann/streamdeck-iconpack-fluentui-system-icons
var elgatoIcons = map[string]iconEntry{
	"com.elgato.streamdeck.multiactions.routine": {
		path:     "./resources/icons/",
		filename: "stack.png",
	},

	"com.elgato.streamdeck.page.goto": {
		path:     "./resources/icons/",
		filename: "",
	},

	"com.elgato.streamdeck.system.website": {
		path:     "./resources/icons/",
		filename: "globe.png",
	},

	"com.elgato.streamdeck.multiactions.routine2": {
		path:     "./resources/icons/",
		filename: "stack_star.png",
	},

	"com.elgato.streamdeck.page.indicator": {
		path:     "./resources/icons/",
		filename: "",
	},

	"com.elgato.streamdeck.system.hotkeyswitch": {
		path:     "./resources/icons/",
		filename: "toggle_right.png",
	},

	"com.elgato.streamdeck.multiactions.random": {
		path:     "./resources/icons/",
		filename: "stack_arrow_forward.png",
	},

	"com.elgato.streamdeck.soundboard.playaudio": {
		path:     "./resources/icons/",
		filename: "play_circle_full.png",
	},

	"com.elgato.streamdeck.system.hotkey": {
		path:     "./resources/icons/",
		filename: "window_new.png",
	},

	"com.elgato.streamdeck.keys.logic": {
		path:     "./resources/icons/",
		filename: "hand_draw.png",
	},

	"com.elgato.streamdeck.soundboard.stopaudioplay": {
		path:     "./resources/icons/",
		filename: "stop.png",
	},

	"com.elgato.streamdeck.system.open": {
		path:     "./resources/icons/",
		filename: "rocket.png",
	},

	"com.elgato.streamdeck.profile.openchild": {
		path:     "./resources/icons/",
		filename: "folder.png",
	},

	"com.elgato.streamdeck.system.timer": {
		path:     "./resources/icons/",
		filename: "clock_alarm.png",
	},

	"com.elgato.streamdeck.system.openapp": {
		path:     "./resources/icons/",
		filename: "app_folder.png",
	},

	"com.elgato.streamdeck.profile.rotate": {
		path:     "./resources/icons/",
		filename: "production.png",
	},

	"com.elgato.streamdeck.system.keybrightness": {
		path:     "./resources/icons/",
		filename: "weather_sunny.png",
	},

	"com.elgato.streamdeck.system.close": {
		path:     "./resources/icons/",
		filename: "dismiss_square.png",
	},

	"com.elgato.streamdeck.page.previous": {
		path:     "./resources/icons/",
		filename: "arrow_circle_left.png",
	},

	"com.elgato.streamdeck.system.sleep": {
		path:     "./resources/icons/",
		filename: "sleep.png",
	},

	"com.elgato.streamdeck.system.text": {
		path:     "./resources/icons/",
		filename: "scan_type.png",
	},

	"com.elgato.streamdeck.page.next": {
		path:     "./resources/icons/",
		filename: "arrow_circle_right.png",
	},

	"com.elgato.streamdeck.system.vsdtoggle": {
		path:     "./resources/icons/",
		filename: "qr_code.png",
	},

	"com.elgato.streamdeck.system.multimedia": {
		path:     "./resources/icons/",
		filename: "play_circle.png",
	},

	"com.elgato.streamdeck.profile.backtoparent": {
		path:     "./resources/icons/",
		filename: "arrow_enter_up.png",
	},
}
