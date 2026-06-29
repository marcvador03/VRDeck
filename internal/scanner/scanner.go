package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type CurrentPage struct {
	Pages struct {
		Current string `json:"Current"`
	} `json:"Pages"`
}

func getCurrentProfile(name string) (string, string, error) {
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

func StartProfilesScan(path string) {
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
					profile, page, err := getCurrentProfile(event.Name)
					if err != nil {
						continue
					}
					log.Info(("New Profile & Page opened"),
						zap.String("profile", profile),
						zap.String("page", page))
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
	err = watcher.Add(path)
	if err != nil {
		log.Error(("Error while adding folder for watch"),
			zap.Error(err))
	}
	<-make(chan struct{})
}
