package utils
import (
	"cmp"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
	"regexp"
	"strconv"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	Start     string    `json:"startDate"`
	End       string    `json:"endDate"`
	StartDate Time
	EndDate   Time
	Location  string    `json:"location"`
	Summary   string    `json:"summary"`

	Name 	  string    `json:"name"`
	Groups    string    `json:"groups"`
	Year      string    `json:"year"`
	DayOfTheWeek string `json:"dayOfTheWeek"`
	Type string 		`json:"type"`
	Parcours string     `json:"parcours"`
}

type Time struct {
	Year    string
	Month   string
	Day     string
	Hour    string
	Minute  string
	Seconde string
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Write(file_name string, dat []byte) {
	err := os.WriteFile(file_name, dat, 0644)
	Check(err)
}
func Read(file_name string) []byte {
	dat, err := os.ReadFile(file_name)
	Check(err)
	return dat
}

func Parse(s []byte) []*Event {
	calendar, err := ics.ParseCalendar(strings.NewReader(string(s)))
	if err != nil {
		fmt.Print("Error in a ics\n")
		return nil
	}
	events := []*Event{}
	for _, e := range calendar.Events() {
		//properties => list of IANAproperties {IANAtoken arg value}
		//filter and create another struct
		event := Event{}
		events = append(events, &event)
		for _, propeties := range e.Properties {
			switch propeties.IANAToken {
			case "DTSTART":
				event.Start = propeties.Value
			case "DTEND":
				event.End = propeties.Value
			case "SUMMARY":
				event.Summary = propeties.Value
			case "LOCATION":
				event.Location = propeties.Value
			}
		}
	}
	return events
}

func Get(code string) []byte {
	resp, err := http.Get(code)
	Check(err)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	return content
}

func Get_time_less_one_day() string {
	now := time.Now()
	// for test
	// now := time.Date(2025, 11, 19, 15, 29, 58, 0, time.Now().Location())
	// now := time.Date(2025, 11, 19, 8, 31, 58, 0, time.Now().Location())

	yesterday := now.Add(-24 * time.Hour) //removve a day
	yesterday2359 := time.Date(
		yesterday.Year(),
		yesterday.Month(),
		yesterday.Day(),
		23, // Heure
		59, // Minute
		0,  // Seconde
		0,  // Nanoseconde
		yesterday.Location(),
	)
	formatted := yesterday2359.Format("20060102T150405")

	return formatted
}

func Sort_events(events []*Event) []*Event {
	if events == nil {return nil}
	slices.SortFunc(events,
		func(e1, e2 *Event) int {
			return cmp.Compare(e1.Start, e2.Start)
		})
	return events
}

func Trunck(events []*Event) []*Event {
	if len(events) == 0 {return nil}
	// we'll use dicothomie search
	time := Get_time_less_one_day()

	var dico_search func(slice []*Event, searched string, start, end int) int
	dico_search = func(slice []*Event, searched string, start, end int) int {
		if end == -1 {
			end = len(slice)
		}

		if end-1 == start {
			return start
		}
		half := (start + end) / 2
		if slice[half].Start < searched {
			return dico_search(slice, searched, half, end)
		} else {
			return dico_search(slice, searched, start, half)
		}
	}

	index := dico_search(events, time, 0, -1)
	events = events[index:]
	return events
}

func ParseDate(date string) Time {
	return Time{date[:4], date[4:6], date[6:8], date[9:11], date[11:13], date[13:15]}
}

func AddDate(events []*Event) []*Event {
	for _, e := range events {
		e.StartDate = ParseDate(e.Start)
		e.EndDate = ParseDate(e.End)

		/* dayOfTheWeek */
		y,err := strconv.Atoi(e.StartDate.Year)
		if err != nil {return nil}
		m,err := strconv.Atoi(e.StartDate.Month)
		if err != nil {return nil}
		month := time.Month(m)
		d,err := strconv.Atoi(e.StartDate.Day)
		if err != nil {return nil}

		e.DayOfTheWeek =
			time.Date(y,month,d,0,0,0,0, time.UTC).Weekday().String()
	}

	
	return events
}

/* depend of the requested ics */
func AddGroups(event *Event, target string){

	event.Groups = target

	/* For the Type */
	upperName := strings.ToUpper(event.Summary)
	if strings.Contains(upperName, "TP") {
		event.Type = "TP"
	} else if strings.Contains(upperName, "CM") {
		event.Type = "CM"
	} else if strings.Contains(upperName, "TD") {
		event.Type = "TD"
	}  else if strings.Contains(upperName, "EXAMEN") {
		event.Type = "EXAMEN"
	}  else if strings.Contains(upperName, "PARTIEL") {
		event.Type = "PARTIEL"
	} else {
		event.Type = "Autre"
	}
}

/* depend of the requested ics */
func AddYear(event *Event, target string){
	upperTarget := strings.ToUpper(target)
	years := []string{"L1", "L2", "L3", "M1", "M2"}
	for _, y := range years {
		if strings.Contains(upperTarget, y) {
			event.Year = y
			return
		}
	}


	if strings.Contains(upperTarget,"S1") || strings.Contains(upperTarget,"S2") {
		event.Year = "L1"
		return
	}
	if strings.Contains(upperTarget,"S3") || strings.Contains(upperTarget,"S4") {
		event.Year = "L2"
		return
	}
	if strings.Contains(upperTarget,"S5") || strings.Contains(upperTarget,"S6") {
		event.Year = "L3"
		return
	}

}

/* depend of the summary */
func AddName(events *Event){
	// remove all TD/TP/CM and number 
	re := regexp.MustCompile(`(?i)(CM|TD|TP)\s?\d{0,2}|L[1-3]|Examen|-`)
	events.Name = re.ReplaceAllString(events.Summary, "")
	r2 := regexp.MustCompile(`\s+`)
	events.Name = r2.ReplaceAllString(events.Name, " ")
	events.Name = strings.TrimSpace(events.Name)
}