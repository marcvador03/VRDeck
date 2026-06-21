package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	l "streamdeckVR/internal/logger"
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

func newProfile(UUID string, path string, data []byte, log *l.Logger) (Profile, error) {
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
		page, err := NewPage(puuid, filepath.Join(path, UUID), log)
		if err != nil {
			continue
		}
		profile.Pages = append(profile.Pages, &page)
	}
	return profile, nil
}

func NewProfileList(path string, log *l.Logger) (*ProfileList, error) {
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
				log.Error(err)
				continue
			}
			if p, err := newProfile(dir.Name(), path, data, log); err != nil {
				log.Error(err)
			} else {
				ProfileList.Profiles = append(ProfileList.Profiles, p)
			}
		}
	}
	return &ProfileList, nil
}

func profileListScan(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)

	return files, err
}
