package ical

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/conf"
	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/mqtt"
	"github.com/apognu/gocal"
)

type EventData struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func fetchAndParseICS(url string, start, end time.Time) (*gocal.Gocal, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ICS: Failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ICS: unexpected status code: %d", resp.StatusCode)
	}

	calendar := gocal.NewParser(resp.Body)
	calendar.Start, calendar.End = &start, &end
	err = calendar.Parse()
	if err != nil {
		return nil, fmt.Errorf("ICS: Failed to parse calendar: %w", err)
	}

	return calendar, nil
}

func getCalendarEvents(url string, name string, start, end time.Time, client *mqtt.MQTTClient) error {
	fmt.Println("ICS: Getting Calendar Events for ", name)
	calendar, err := fetchAndParseICS(url, start, end)
	if err != nil {
		return fmt.Errorf("ICS: Failed to fetch and parse ICS: %w", err)
	}

	fmt.Println("ICS: Found ", len(calendar.Events), " events in ", name)
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1)
	fmt.Println("ICS: Today Start: ", todayStart)
	fmt.Println("ICS: Today End: ", todayEnd)

	var todayEvents []gocal.Event
	for _, event := range calendar.Events {
		if event.Start.After(todayStart) && event.Start.Before(todayEnd) {
			todayEvents = append(todayEvents, event)
		}
	}

	fmt.Println("ICS: Events for today:")
	for _, event := range todayEvents {
		fmt.Printf("ICS: %s on %s Location: %s", event.Summary, event.Start, event.Location)
		fmt.Println()
	}
	fmt.Println()

	if len(todayEvents) == 0 {
		topic := fmt.Sprintf(name + "/today/events")
		fmt.Println("ICS: No events for today")
		mqtt.Publish(*client, topic, "[]", false)
		return nil
	} else {
		var eventsData []EventData
		for _, event := range todayEvents {
			eventsData = append(eventsData, EventData{
				Summary: event.Summary,
				Start:   event.Start.Format("15:04"),
				End:     event.End.Format("15:04"),
			})
		}
		eventsJSON, err := json.Marshal(eventsData)
		if err != nil {
			return fmt.Errorf("ICS: Failed to marshal event data: %w", err)
		}
		topic := fmt.Sprintf(name + "/today/events")
		mqtt.Publish(*client, topic, string(eventsJSON), false)
	}

	return nil

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

	fmt.Println("ICS: Setting up ICS Events")
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
					getCalendarEvents(ics.URL, ics.Name, start, end, client)
				}
			}
		}()
	}
}
