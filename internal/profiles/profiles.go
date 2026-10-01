package profiles

import (
	"VRDeck/internal/logger"
	"VRDeck/internal/ws"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

func NewProfileList(path string, ws *ws.MSFSWebSocket) *ProfileList {
	return &ProfileList{
		path: path,
		ws:   ws,
	}
}

func (p *ProfileList) getPageByUUID(profileUUID, pageUUID string) (*Pages, *Pages) {
	log := logger.GetDefaultLogger()
	for _, profile := range p.Profiles {
		if profile.UUID == profileUUID {
			for _, page := range profile.Pages {
				if page.UUID == pageUUID {
					return page, profile.Default
				}
			}
			log.Error(("Page not found in Profile"),
				zap.String("profile", profile.Name),
				zap.String("page", pageUUID))
			return nil, nil
		}
	}
	log.Error(("Profile not found"),
		zap.String("profile", profileUUID))
	return nil, nil
}

func (p *ProfileList) getProfileByUUID(profileUUID string) *Profile {
	log := logger.GetDefaultLogger()
	for _, profile := range p.Profiles {
		if strings.TrimSuffix(profile.UUID, ".sdProfile") == profileUUID {
			return &profile
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
			Default string   `json:"Default"`
			Current string   `json:"Current"`
			Pages   []string `json:"Pages"`
		} `json:"Pages"`
	}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return profile, err
	}
	profile.UUID = strings.TrimSuffix(UUID, ".sdProfile")
	profile.Name = rawData.Name
	profile.Default = nil
	profile.Current = nil
	profile.pagesNum = len(rawData.Pages.Pages)
	for _, puuid := range rawData.Pages.Pages {
		page, err := p.newPage(puuid, filepath.Join(p.path, UUID))
		if err != nil {
			continue
		}
		profile.Pages = append(profile.Pages, &page)
		if rawData.Pages.Current == puuid {
			profile.Current = &page
		}
	}
	page, err := p.newPage(rawData.Pages.Default, filepath.Join(p.path, UUID))
	if err == nil {
		profile.Pages = append(profile.Pages, &page)
		profile.Default = &page
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
	var latest time.Time = time.Unix(0, 0)
	p.Current = ""
	for _, dir := range profiledir {
		if dir.IsDir() {
			fs, err := dir.Info()
			if err != nil {
				log.Error("Error while reading directory info",
					zap.String("UUID", dir.Name()),
					zap.Error(err))
				continue
			}
			if fs.ModTime().After(latest) {
				latest = fs.ModTime()
				p.Current = dir.Name()
			}
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
				log.Info("Added a New Profile to the list: ",
					zap.String("Name", profile.Name),
					zap.Int("Number of pages", len(profile.Pages)))
			}
		}
	}
	log.Info("Number of Profiles scanned and stored: ",
		zap.Int("Number of pages", len(p.Profiles)))
	return nil
}
