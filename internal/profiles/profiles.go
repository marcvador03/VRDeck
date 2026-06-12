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
	name     string
	pagesNum int
	page     []Pages
}

type Pages struct {
	profile *Profile
	UUID    string
}

type ProfileList struct {
	Profiles []Profile
}

func newProfile(UUID, name string) Profile {
	return Profile{
		UUID: UUID,
		name: name,
	}
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
			var manifest map[string]interface{}
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				log.Error(err)
				continue
			}
			if err := json.Unmarshal(data, &manifest); err != nil {
				log.Error(err)
			}
			if name, ok := manifest["Name"].(string); ok {
				p := newProfile(strings.TrimSuffix(dir.Name(), "."), name)
				ProfileList.Profiles = append(ProfileList.Profiles, p)
			} else {
				log.Error(fmt.Errorf("Field 'name' not found or not a string"))
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
