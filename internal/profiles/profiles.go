package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	l "streamdeckVR/internal/logger"
	"strings"
)

type Profile struct {
	UUID     string
	Name     string
	pagesNum int
	Pages    []*Pages
}

type Pages struct {
	UUID string
}

type ProfileList struct {
	Profiles []Profile
}

func newProfile(UUID string, data []byte) (Profile, error) {
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
		page := &Pages{
			UUID: puuid,
		}
		profile.Pages = append(profile.Pages, page)
	}
	return profile, nil
}

func NewProfileList(path string, log *l.Logger) (*ProfileList, error) {
	var ProfileList ProfileList
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
			if p, err := newProfile(strings.TrimSuffix(dir.Name(), "."), data); err != nil {
				log.Error(err)
			} else {
				ProfileList.Profiles = append(ProfileList.Profiles, p)
			}
		}
		// fmt.Printf(
		// 	"Name: %s,  : %v, Type: %v\n",
		// 	dir.Name(),
		// 	dir.IsDir(),
		// 	dir.Type(),
		// )
		// fmt.Println(dir.Name())
	}
	return &ProfileList, nil
}

func profileListScan(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)

	return files, err
}
