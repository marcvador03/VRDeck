package profiles

import (
	"VRDeck/internal/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func (p *ProfileList) mergeDefaut(page *Pages, defaut *Pages) Pages {
	tmp := *page
	if defaut.Buttons == nil {
		return tmp
	}
	tmp.Buttons = nil
	for _, defbutton := range defaut.Buttons {
		tmp.Buttons = append(tmp.Buttons, defbutton)
	}
	for _, button := range page.Buttons {
		for _, defbutton := range defaut.Buttons {
			if defbutton.Col == button.Col && defbutton.Row == button.Row {
				continue
			} else {
				tmp.Buttons = append(tmp.Buttons, button)
			}
		}
	}
	return tmp
}

func (p *ProfileList) sendPageUpdate(page *Pages) {
	log := logger.GetDefaultLogger()

	jsonData, err := json.Marshal(page)
	if err != nil {
		log.Error(("failed to marshal page to JSON"),
			zap.Error(err))
		return
	}
	if err := p.ws.BroadcastJSON(jsonData); err != nil {
		log.Error(("failed to broadcast JSON via WebSocket"),
			zap.Error(err))
		return
	}
}

func (p *ProfileList) readCurrentManifest(manifestPath string, profileUUID string) ([]byte, error) {
	log := logger.GetDefaultLogger()
	for i := 0; i < 10; i++ {
		data, err := os.ReadFile(manifestPath)
		if err == nil {
			return data, nil
		}
		if os.IsNotExist(err) || strings.Contains(err.Error(), "being used by another process") {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		log.Error("Error while reading manifest file",
			zap.String("UUID", profileUUID),
			zap.Error(err))
		return nil, err
	}
	log.Error("Error while reading manifest file after 10 trials",
		zap.String("UUID", profileUUID))
	return nil, fmt.Errorf("Error while reading manifest file after 10 trials")
}

func (p *ProfileList) getCurrentProfile(name string) (string, string, error) {
	log := logger.GetDefaultLogger()
	profileUUID := strings.TrimSuffix(filepath.Base(name), ".sdProfile")
	var page string
	data, err := p.readCurrentManifest(filepath.Join(name, "manifest.json"), profileUUID)
	if err != nil {
		return profileUUID, page, err
	}
	var current CurrentPage
	err = json.Unmarshal([]byte(data), &current)
	if err != nil {
		log.Error(("Error while reading Json"),
			zap.String("UUID", profileUUID),
			zap.Error(err))
		return profileUUID, page, err
	}
	page = current.Pages.Current
	return profileUUID, page, nil
}

func (p *ProfileList) StartProfilesScan() {
	log := logger.GetDefaultLogger()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(("Error while starting FsNotify watcher"),
			zap.Error(err))
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					p.GetStartingPage(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(("Error while watching folder"),
					zap.Error(err))
			}
		}
	}()
	err = watcher.Add(p.path)
	if err != nil {
		log.Error(("Error while adding folder for watch"),
			zap.Error(err))
	}
	<-make(chan struct{})
}

func (p *ProfileList) GetStartingPage(name string) error {
	log := logger.GetDefaultLogger()
	var page, defaut *Pages
	profileUUID, pageUUID, err := p.getCurrentProfile(name)
	if err != nil {
		log.Info("Error while opening Profile, defaulting to current values stored")
		currentPage := p.getProfileByUUID(profileUUID).Current
		page, defaut = p.getPageByUUID(profileUUID, currentPage.UUID)

		log.Info(("New Profile & Page opened"),
			zap.String("profile", profileUUID),
			zap.String("page", currentPage.UUID))
	} else {
		page, defaut = p.getPageByUUID(profileUUID, pageUUID)
		log.Info(("New Profile & Page opened"),
			zap.String("profile", profileUUID),
			zap.String("page", pageUUID))
	}
	if page == nil {
		log.Error(("Error while opening Page"),
			zap.Error(err))
		return err
	}
	pageToSend := p.mergeDefaut(page, defaut)
	p.sendPageUpdate(&pageToSend)
	return nil
}
