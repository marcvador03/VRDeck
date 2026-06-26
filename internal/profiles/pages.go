package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"streamdeckVR/internal/logger"
	"strings"
	"unicode"

	"go.uber.org/zap"
)

type Buttons struct {
	Row      int
	Col      int
	Title    string
	ActionID string
}

type Pages struct {
	UUID    string
	Name    string
	Buttons []Buttons
}

type Controller struct {
	Buttons map[string]map[string]interface{} `json:"Actions"`
}

func convertTiletoDigits(tile string) (int, int, error) {
	rStr, cStr, ok := strings.Cut(tile, ",")
	fmt.Printf("r:" + rStr + "c: " + cStr)
	if !ok || len(rStr) != 1 || len(cStr) != 1 {
		return -1, -1, fmt.Errorf("Error 1")
	}
	if !unicode.IsDigit(rune(rStr[0])) || !unicode.IsDigit(rune(cStr[0])) {
		return -1, -1, fmt.Errorf("Error 2")
	}
	r, c := int(rStr[0]-'0'), int(cStr[0]-'0')
	return r, c, nil
}

func addButton(tile string, details map[string]interface{}) (Buttons, error) {
	Button := Buttons{}
	row, col, err := convertTiletoDigits(tile)
	if err != nil {
		return Buttons{}, err
	}
	Button.Row, Button.Col = row, col
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
		Button.Title = "unk"
		return Button, nil
		//return Button, fmt.Errorf("Title is missing or not a string")
	}
	Button.Title = title
	return Button, nil
}

func addDetailPage(page *Pages, data []byte) error {
	var rawData struct {
		Controllers []Controller `json:"Controllers"`
		Name        string       `json:"Name"`
	}

	if err := json.Unmarshal(data, &rawData); err != nil {
		return fmt.Errorf("Error while unpacking json for Actions %w", err)
	}
	if len(rawData.Controllers) > 0 {
		page.Name = rawData.Name
		for tile, details := range rawData.Controllers[0].Buttons {
			NewAction, err := addButton(tile, details)
			if err != nil {
				return fmt.Errorf("Wrong format in Actions %s %w", page.UUID, err)
			}
			page.Buttons = append(page.Buttons, NewAction)
		}
	}
	return nil
}

func NewPage(UUID string, path string) (Pages, error) {
	log := logger.GetDefaultLogger()
	page := Pages{
		UUID: UUID,
	}
	pagePath := filepath.Join(path, "Profiles", UUID)
	pagedir, err := os.ReadDir(pagePath)
	if err != nil {
		return Pages{}, fmt.Errorf("failed to read path %s: %w", pagePath, err)
	}
	for _, dir := range pagedir {
		if dir.IsDir() {
			manifestPath := filepath.Join(pagePath, "manifest.json")
			//fmt.Printf(manifestPath)
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				log.Error("Error while reading manifest file",
					zap.String("path", pagePath),
					zap.Error(err))
				continue
			}
			if err := addDetailPage(&page, data); err != nil {
				log.Error("failed to add details",
					zap.String("page", dir.Name()),
					zap.Error(err))
			}
		}
	}
	return page, nil
}
