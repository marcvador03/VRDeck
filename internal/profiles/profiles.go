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

func modifyProfileManifest(manifestPath string, data []byte, name string) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if _, ok := raw["InstalledByPluginUUID"].(string); !ok {
		raw["InstalledByPluginUUID"] = "com.streamdeckVR.plugin"
		if _, ok := raw["PreconfiguredName"].(string); !ok {
			raw["PreconfiguredName"] = name
		}
		if _, ok := raw["ReadOnly"].(string); !ok {
			raw["ReadOnly"] = false
		}
	} else {
		return fmt.Errorf("Profile is already owned by a plugin %s", raw["InstalledByPluginUUID"])
	}
	newJson, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	if err := os.WriteFile(manifestPath, newJson, 0644); err != nil {
		return fmt.Errorf("Error while updating Profile manifest: %w", err)
	}
	return nil
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
				if err := modifyProfileManifest(manifestPath, data, p.Name); err != nil {
					log.Error(err)
				}
			}
		}
	}
	return &ProfileList, nil
}
