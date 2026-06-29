package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"
	"streamdeckVR/internal/ws"

	"go.uber.org/zap"
)

func NewProfileList(path string, ws *ws.MSFSWebSocket) *ProfileList {
	return &ProfileList{
		path: path,
		ws:   ws,
	}
}

func (p *ProfileList) getPageByUUID(profileUUID, pageUUID string) *Pages {
	log := logger.GetDefaultLogger()
	for _, profile := range p.Profiles {
		if profile.UUID == profileUUID {
			for _, page := range profile.Pages {
				if page.UUID == pageUUID {
					return page
				}
			}
			log.Error(("Page not found in Profile"),
				zap.String("profile", profile.Name),
				zap.String("page", pageUUID))
			return nil
		}
	}
	log.Error(("Profile not found"),
		zap.String("profile", profileUUID))
	return nil
}

func (p *ProfileList) InspectData() {
	log := logger.GetDefaultLogger()
	for i, profile := range p.Profiles {
		log.Info(("Profile"),
			zap.Int("Index", i),
			zap.String("Profile Name", profile.Name),
			zap.String("Profile UUID", profile.UUID))
		for j, page := range profile.Pages {
			log.Info(("Page"),
				zap.Int("Index", j),
				zap.String("Page Name", page.Name),
				zap.String("Page UUID", page.UUID))
			for k, button := range page.Buttons {
				log.Info(("Buttons"),
					zap.Int("Index", k),
					zap.Int("Row", button.Row),
					zap.Int("Col", button.Col),
					zap.String("Title", button.Title),
					zap.String("ActionID", button.ActionID))
			}
		}
	}
}

func (p *ProfileList) newProfile(UUID string, data []byte) (Profile, error) {
	var profile Profile
	var rawData struct {
		Name  string `json:"Name"`
		Pages struct {
			Pages []string `json:"Pages"`
		} `json:"Pages"`
	}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return profile, err
	}
	profile.UUID = UUID
	profile.Name = rawData.Name
	profile.pagesNum = len(rawData.Pages.Pages)
	for _, puuid := range rawData.Pages.Pages {
		page, err := p.newPage(puuid, filepath.Join(p.path, UUID))
		if err != nil {
			continue
		}
		profile.Pages = append(profile.Pages, &page)
	}
	return profile, nil
}

func (p *ProfileList) CreateProfileList() error {
	log := logger.GetDefaultLogger()
	profiledir, err := os.ReadDir(p.path)
	if err != nil {
		log.Error(("failed to read path"),
			zap.String("path", p.path),
			zap.Error(err))
		return err
	}

	for _, dir := range profiledir {
		if dir.IsDir() {
			manifestPath := p.path + "\\" + dir.Name() + "\\" + "manifest.json"
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				log.Error("Error while reading manifest file",
					zap.String("UUID", dir.Name()),
					zap.Error(err))
				continue
			}
			if profile, err := p.newProfile(dir.Name(), data); err != nil {
				log.Error("Error while process profile file",
					zap.String("UUID", dir.Name()),
					zap.Error(err))
			} else {
				p.Profiles = append(p.Profiles, profile)
			}
		}
	}
	return nil
}
