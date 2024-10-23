package tools

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/mqtt"
	log "github.com/sirupsen/logrus"
)

type EventData struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

type CalendarEvent struct {
	Summary    string
	Start      time.Time
	End        time.Time
	Reacurring bool
	UID        string
	TimeZone   string
}

func eventsToJSON(events []CalendarEvent) (string, error) {
	var eventsData []EventData
	for _, event := range events {
		eventsData = append(eventsData, EventData{
			Summary: event.Summary,
			Start:   event.Start.Format("15:04"),
			End:     event.End.Format("15:04"),
		})
	}
	eventsJSON, err := json.Marshal(eventsData)
	if err != nil {
		return "[]", fmt.Errorf("cal tools: failed to marshal event data: %w", err)
	}
	return string(eventsJSON), nil
}

func filterEventsToday(events []CalendarEvent) []CalendarEvent {
	var todayEvents []CalendarEvent
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1)
	//fmt.Println("Cal Tools: Today Start: ", todayStart)
	//fmt.Println("Cal Tools: Today End: ", todayEnd)
	for _, event := range events {
		if event.Start.After(todayStart) && event.Start.Before(todayEnd) {
			todayEvents = append(todayEvents, event)
			fmt.Printf("Cal Tools: Added: %s \n Start Time: %s \n Reaccuring: %t \n UID: %s \n", event.Summary, event.Start, event.Reacurring, event.UID)
		}
	}
	return todayEvents
}

func filterOutdatedTodayEvents(events []CalendarEvent) []CalendarEvent {
	// This function expects the input to be events for today
	var todayFilteredEvents []CalendarEvent
	now := time.Now()
	//fmt.Println("Cal Tools: Now: ", location)
	for _, event := range events {
		location, err := time.LoadLocation(event.TimeZone)
		if err != nil {
			log.Error("Cal Tools: Error loading location: ", err)
			location = time.Local
		}
		eventEndLocal := event.End.In(location)
		//fmt.Println("Event TZ: ", event.TimeZone)
		if eventEndLocal.After(now) {
			todayFilteredEvents = append(todayFilteredEvents, event)
			//fmt.Printf("Cal Tools: Added: %s \n UID: %s \n Start: %s \n End: %s \n Reacurring: %t \n Time: %s\n", event.Summary, event.UID, event.Start, event.End, event.Reacurring, now)
		} /* else {
			fmt.Printf("Cal Tools: Skipped: %s Start: %s End: %s Time: %s\n", event.Summary, event.Start, event.End, now)
		}
		*/

	}
	return todayFilteredEvents
}

func getErrorJSON() string {
	eventsData := []EventData{
		{
			Summary: "Error",
			Start:   "00:00",
			End:     "23:59",
		},
	}
	eventsJSON, err := json.Marshal(eventsData)
	if err != nil {
		// Dont see why this would ever error #famouslastwords
		log.Error("cal tools: failed to marshal event data: %w", err)
		return "[Server Error]"
	}
	return string(eventsJSON)
}

func PublishCalendarEvents(client mqtt.MQTTClient, name string, events []CalendarEvent) {
	// Just pipe all the fetched events into this function and it will filter and publish the appropriate events

	// Define the topics
	var upcoming string = fmt.Sprintf("%s/today/upcoming/events", name)
	var all string = fmt.Sprintf("%s/today/all/events", name)

	// Map of topics to publish
	topics := map[string]string{
		all:      "",
		upcoming: "",
	}

	// Get all todays events
	todayEvents := filterEventsToday(events)
	if len(todayEvents) == 0 {
		fmt.Println("cal tools: No events for today for ", name)
		topics[all] = "[]"
		topics[upcoming] = "[]"
	} else {
		todayJSON, err := eventsToJSON(todayEvents)
		if err != nil {
			log.Error("Cal Tools: Error publishing all day events for: ", name, " Error: ", err)
			topics[all] = getErrorJSON()
		} else {
			topics[all] = todayJSON
		}
		upcomingJSON, err := eventsToJSON(filterOutdatedTodayEvents(todayEvents))
		if err != nil {
			log.Error("Cal Tools: Error publishing upcoming events for: ", name, " Error: ", err)
			topics[upcoming] = getErrorJSON()
		} else {
			topics[upcoming] = upcomingJSON
		}
	}

	// Loop through the map and publish the topics
	for topic, value := range topics {
		mqtt.Publish(client, topic, value, false)
	}
	// add new line for readability
	println()
}
