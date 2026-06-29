package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func (p *ProfileList) sendPageUpdate(page *Pages) {
	log := logger.GetDefaultLogger()
	jsonData, err := json.Marshal(page)
	if err != nil {
		log.Error(("failed to marshal page to JSON"),
			zap.Error(err))
		return
	}

	// Broadcast the JSON to all connected clients
	if err := p.ws.BroadcastJSON(jsonData); err != nil {
		log.Error(("failed to broadcast JSON via WebSocket"),
			zap.Error(err))
		return
	}
}

func (p *ProfileList) getCurrentProfile(name string) (string, string, error) {
	log := logger.GetDefaultLogger()
	manifestPath := filepath.Join(name, "manifest.json")
	profile := strings.TrimSuffix(filepath.Base(name), ".sdProfile")
	var page string
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Error("Error while reading manifest file",
			zap.String("UUID", profile),
			zap.Error(err))
		return profile, page, err
	}
	var current CurrentPage
	err = json.Unmarshal([]byte(data), &current)
	if err != nil {
		log.Error(("Error while reading Json"),
			zap.String("UUID", profile),
			zap.Error(err))
		return profile, page, err
	}
	page = current.Pages.Current
	return profile, page, nil
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
					profileUUID, pageUUID, err := p.getCurrentProfile(event.Name)
					if err != nil {
						continue
					}
					page := p.getPageByUUID(profileUUID, pageUUID)
					if page == nil {
						continue
					}
					log.Info(("New Profile & Page opened"),
						zap.String("profile", profileUUID),
						zap.String("page", pageUUID))
					p.sendPageUpdate(page)
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
