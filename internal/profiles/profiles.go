package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"

	"go.uber.org/zap"
)

type Profile struct {
	UUID     string
	Name     string
	pagesNum int
	Pages    []*Pages
}

type ProfileList struct {
	path     string
	Profiles []Profile
}

func newProfile(UUID string, path string, data []byte) (Profile, error) {
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
		page, err := NewPage(puuid, filepath.Join(path, UUID))
		if err != nil {
			continue
		}
		profile.Pages = append(profile.Pages, &page)
	}
	return profile, nil
}

func NewProfileList(path string) (*ProfileList, error) {
	log := logger.GetDefaultLogger()
	ProfileList := ProfileList{
		path: path,
	}
	profilev3, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path %s: %w", path, err)
	}

	for _, dir := range profilev3 {
		if dir.IsDir() {
			manifestPath := path + "\\" + dir.Name() + "\\" + "manifest.json"
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				log.Error("Error while reading manifest file",
					zap.String("UUID", dir.Name()),
					zap.Error(err))
				continue
			}
			if p, err := newProfile(dir.Name(), path, data); err != nil {
				log.Error("Error while process profile file",
					zap.String("UUID", dir.Name()),
					zap.Error(err))
			} else {
				ProfileList.Profiles = append(ProfileList.Profiles, p)
			}
		}
	}
	return &ProfileList, nil
}
