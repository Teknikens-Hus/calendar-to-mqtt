package ical

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/calendars/tools"
	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/conf"
	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/mqtt"
	"github.com/apognu/gocal"
	log "github.com/sirupsen/logrus"
)

func fetchAndParseICS(url string, start, end time.Time) (*gocal.Gocal, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Check if the content type is text/calendar, if its text/html the URL is probably wrong or expired
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/calendar; charset=utf-8" {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}

	calendar := gocal.NewParser(resp.Body)

	// Here we can map the timezone IDs from the ICS file to the Go time.Location
	// This is useful if yourtimezone cant be resolved by Go
	var tzMapping = map[string]string{
		"W. Europe Standard Time": "Europe/Stockholm",
	}
	gocal.SetTZMapper(func(s string) (*time.Location, error) {
		if tzid, ok := tzMapping[s]; ok {
			return time.LoadLocation(tzid)
		}
		return nil, fmt.Errorf("")
	})
	// Set the start and end date for the calendar parser (Which event dates to parse)
	calendar.Start, calendar.End = &start, &end
	err = calendar.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar: %w", err)
	}
	return calendar, nil
}

func getICSEvents(url string, start, end time.Time) ([]tools.CalendarEvent, error) {
	calendar, err := fetchAndParseICS(url, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch and parse ICS: %w", err)
	}
	// Convert the gocal events to our own CalendarEvent struct
	var events []tools.CalendarEvent
	for _, event := range calendar.Events {
		events = append(events, tools.CalendarEvent{
			Summary:    event.Summary,
			Start:      *event.Start,
			End:        *event.End,
			Reacurring: event.IsRecurring,
			UID:        event.Uid,
		})

		if event.IsRecurring {
			fmt.Println("ICS: Recurring event: ", event.Summary)
			fmt.Println("ICS: Recurring exclude date num: ", len(event.ExcludeDates))
			/*
				for _, excludedDate := range event.ExcludeDates {
					fmt.Println("ICS: Excluded date: ", excludedDate)
				}
				fmt.Printf("ICS: Recurring rule: %s \n %s \n UID: %s", event.RecurrenceID, event.RecurrenceRule, event.Uid)
			*/
		}

	}
	return events, nil
}

func SetupICS(client *mqtt.MQTTClient) {
	icsConf, err := conf.GetICSConfig()
	if err != nil {
		fmt.Println("ICS: Error setting up ICS: ", err)
		return
	}

	if icsConf == nil {
		fmt.Println("ICS: No ICS configurations found in the config file")
		return
	}

	// Calculate the first and last day of the current month
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	start := firstDay
	end := lastDay
	fmt.Println("ICS: Start: ", start)
	fmt.Println("ICS: End: ", end)

	fmt.Printf("ICS: Setting up ICS %d Calendar(s) \n", len(*icsConf))
	for _, ics := range *icsConf {
		// Create a new ticker for each ICS
		fmt.Println("ICS: Setting up ticker for ", ics.Name, " with interval ", ics.Interval, " seconds")
		ticker := time.NewTicker(time.Duration(ics.Interval) * time.Second)
		//defer ticker.Stop()

		// Run scheduled task in a goroutine
		go func() {
			for {
				select {
				case <-ticker.C:
					//fmt.Println("ICS: Ticker ticked at ", t)
					events, err := getICSEvents(ics.URL, start, end)
					if err != nil {
						log.Error("ICS: Error getting ICS events: ", err)
					} else {
						fmt.Printf("ICS: Found %d events in %s for date %s to %s \n", len(events), ics.Name, start.Format("2006-01-02"), end.Format("2006-01-02"))
						tools.PublishCalendarEvents(*client, ics.Name, events)
					}
				}
			}
		}()
	}
}
