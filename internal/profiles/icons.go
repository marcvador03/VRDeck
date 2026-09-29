package profiles

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func (p *ProfileList) getImg(imageUrl string, pagePath string) string {
	path := filepath.Join(pagePath, imageUrl)
	img, err := os.ReadFile(path)
	if err != nil {
		fmt.Print("error redaing file")
		return ""
	}
	mime := http.DetectContentType(img)
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(img)
}
