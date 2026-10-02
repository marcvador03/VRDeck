package profiles

import (
	"VRDeck/internal/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"go.uber.org/zap"
)

type Controller struct {
	Buttons map[string]map[string]interface{} `json:"Actions"`
}

func (p *Profile) convertTiletoDigits(tile string) (int, int, error) {
	cStr, rStr, ok := strings.Cut(tile, ",")
	if !ok || len(rStr) != 1 || len(cStr) != 1 {
		return -1, -1, fmt.Errorf("Error 1")
	}
	if !unicode.IsDigit(rune(rStr[0])) || !unicode.IsDigit(rune(cStr[0])) {
		return -1, -1, fmt.Errorf("Error 2")
	}
	r, c := int(rStr[0]-'0'), int(cStr[0]-'0')
	return r, c, nil
}

func (p *Profile) getSettings(Button *Buttons, details map[string]interface{}, pagePath string) {
	settings, ok := details["Settings"].(map[string]interface{})
	if !ok {
		return
	}
	childUUID, ok := settings["ProfileUUID"].(string)
	if ok {
		parentPath := filepath.Dir(filepath.Dir(pagePath))
		childPage, err := p.newPage(childUUID, parentPath)
		if err == nil {
			p.Pages = append(p.Pages, &childPage)
			Button.ChildPage = &childPage
		}
	}
	pageIndex, ok := settings["PageIndex"].(float64)
	if ok {
		Button.PageIndex = strconv.Itoa(int(pageIndex))
	}
}

func (p *Profile) addButton(tile string, details map[string]interface{}, pagePath string) (Buttons, error) {
	Button := Buttons{}
	row, col, err := p.convertTiletoDigits(tile)
	if err != nil {
		return Buttons{}, err
	}
	Button.Row, Button.Col = row, col
	if UUID, ok := details["UUID"].(string); !ok {
		return Buttons{}, fmt.Errorf("Issue with UUID of Action")
	} else {
		Button.UUID = UUID
	}
	if actionID, ok := details["ActionID"].(string); !ok {
		return Buttons{}, fmt.Errorf("Issue with ActionID")
	} else {
		Button.ActionID = actionID
	}
	states, ok := details["States"].([]interface{})
	if !ok || len(states) == 0 {
		return Buttons{}, fmt.Errorf("States is missing, empty, or not an array")
	}
	firstState, ok := states[0].(map[string]interface{})
	if !ok {
		return Buttons{}, fmt.Errorf("First state is not an object")
	}
	title, ok := firstState["Title"].(string)
	if !ok {
		Button.Title = ""
	} else {
		Button.Title = title
	}
	p.getSettings(&Button, details, pagePath)
	p.getImage(&Button, firstState, pagePath)
	return Button, nil
}

func (p *Profile) addDetailPage(page *Pages, data []byte) error {
	var rawData struct {
		Controllers []Controller `json:"Controllers"`
		Name        string       `json:"Name"`
	}

	if err := json.Unmarshal(data, &rawData); err != nil {
		return fmt.Errorf("Error while unpacking json for Actions %w", err)
	}
	//Actions in Streamdecks jsons are Buttons here
	if len(rawData.Controllers) > 0 {
		page.Name = rawData.Name
		for tile, details := range rawData.Controllers[0].Buttons {
			NewButton, err := p.addButton(tile, details, page.pagePath)
			if err != nil {
				return fmt.Errorf("Wrong format in Actions %s %w", page.UUID, err)
			}
			page.Buttons = append(page.Buttons, NewButton)
		}
	}
	return nil
}

func (p *Profile) newPage(UUID string, path string) (Pages, error) {
	log := logger.GetDefaultLogger()
	page := Pages{
		UUID:     UUID,
		pagePath: filepath.Join(path, "Profiles", UUID),
	}
	pagedir, err := os.ReadDir(page.pagePath)
	if err != nil {
		return Pages{}, fmt.Errorf("failed to read path %s: %w", page.pagePath, err)
	}
	for _, dir := range pagedir {
		if dir.IsDir() {
			manifestPath := filepath.Join(page.pagePath, "manifest.json")
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				log.Error("Error while reading manifest file",
					zap.String("path", page.pagePath),
					zap.Error(err))
				continue
			}
			if err := p.addDetailPage(&page, data); err != nil {
				log.Error("failed to add details",
					zap.String("page", dir.Name()),
					zap.Error(err))
			}
		}
	}
	return page, nil
}
