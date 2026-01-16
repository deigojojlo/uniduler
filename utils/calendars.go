package utils
import (
	"os"
	"encoding/json"
)

type CalendarsEntry struct {
	Code     string   `json:"code"`
	Title    string   `json:"title"`
	Label    string   `json:"label"`
	Parcours string   `json:"parcours"`
	YearRaw  string   `json:"year"`
	Page     []string `json:"page"`
}

func ReadCalendars() []CalendarsEntry {
	file, err := os.ReadFile("calendars.json")
	if err != nil {
		panic(err)
	}

	var raw []CalendarsEntry
	err = json.Unmarshal(file, &raw)
	if err != nil {
		panic(err)
	}
	return raw
}