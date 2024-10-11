package ical

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type CalendarEvent struct {
	UID      string
	SUMMARY  string
	DTSTART  time.Time
	DTEND    time.Time
	RRULE    string
	CLASS    string
	PRIORITY int
	DTSTAMP  string // Maybe time?  20241011T130109Z
	TRANSP   string // OPAQUE etc
	STATUS   string // CONFIRMED etc
	SEQUENCE int    // 0
	LOCATION string
}

func GetCalendarEvents(url string) {
	fmt.Println("Getting calendar events from: ", url)

	cphLoc, locErr := time.LoadLocation("Europe/Stockholm")
	if locErr != nil {
		log.Fatalf("Error loading timezone: %v", locErr)
	}
	fmt.Println("Timezone: ", cphLoc)

	events := parseCalendar(fetchCalendar(url))
	fmt.Println("Event 4: ", events[4])

}

func fetchCalendar(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting calendar events: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error getting calendar events: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	return string(body)
}

func parseCalendar(body string) []string {
	fmt.Println("Parsing calendar")

	var events []string

	// Split the bodyt into individual events
	eventStrings := strings.Split(body, "BEGIN:VEVENT")
	for _, eventString := range eventStrings[1:] {
		eventString = "BEGIN:VEVENT" + eventString
		eventString = strings.Split(eventString, "END:VEVENT")[0] + "END:VEVENT"

		//event := CalendarEvent{
		//	UID:      getField(eventString, "UID"),
		//	SUMMARY:  getField(eventString, "SUMMARY"),
		//	DTSTART:  parseTime(getField(eventString, "DTSTART")),
		//	DTEND:    parseTime(getField(eventString, "DTEND")),
		//	CLASS:    getField(eventString, "CLASS"),
		//	PRIORITY: getField(eventString, "PRIORITY"),
		//	DTSTAMP:  getField(eventString, "DTSTAMP"),
		//	TRANSP:   getField(eventString, "TRANSP"),
		//	STATUS:   getField(eventString, "STATUS"),
		//	SEQUENCE: getField(eventString, "SEQUENCE"),
		//	LOCATION: getField(eventString, "LOCATION"),
		//}

		events = append(events, eventString)
	}

	return events
}
