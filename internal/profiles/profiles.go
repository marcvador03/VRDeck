package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	l "streamdeckVR/internal/logger"
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
	profiles []Profile
}

func NewProfileList(path string, log *l.Logger) (*ProfileList, error) {
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
			// Extract the "Name" field
			if name, ok := manifest["Name"].(string); ok {
				fmt.Printf("Name: %s\n", name)
			} else {
				fmt.Println("Field 'name' not found or not a string")
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
	return &ProfileList{}, nil
}

func profileListScan(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)

	return files, err
}
