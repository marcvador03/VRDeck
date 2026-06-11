package profiles

import (
	"os"
)

type Profile struct {
	UUID  string
	name  string
	pages int
}

func ProfileScan(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)

	return files, err
}
