package scanner

import (
	"fmt"
	"streamdeckVR/internal/logger"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

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
				fmt.Println("event:", event)
				if event.Has(fsnotify.Write) {
					fmt.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path)
	if err != nil {
		log.Error(("Error while adding folder for watch"),
			zap.Error(err))
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
