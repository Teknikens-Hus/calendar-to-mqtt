// calendar-to-mqtt is a simple application that reads a Calendars (ICS etc) and publishes the events to an MQTT broker.
package main

import (
	"fmt"

	"os"
	"time"

	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/mqtt"
	"github.com/apognu/gocal"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetReportCaller(true)
	// Run the application
	fmt.Println("Starting calendar-to-mqtt...")
	fmt.Println("Created by:")
	fmt.Println(`                                                                                                      
 #############    ########     ###      #### ############    ####  ###      ####   ########     ############     ########       ####       ###  ###       ####   ########  
############## #############   ###    ###### ##############  ####  ###    ####################  #############  ############     ####       ###  ###       #### ############
      ###     ######    #####  ###  ######   #####     ####  ####  ###  ######  #####    ###### ####     ##### ####    ####     ####       ###  ###       #### ####    ####
      ###     ################ #########     ####       ###  ####  #########   ################ ###       #### ########         ##############  ###       #### ########    
      ###     ################ #########     ####       ###  ####  #########   ################ ###       ####  ##########      ##############  ###       ####  ########## 
      ###     ####        ##   ##### #####   ####       ###  ####  ##### ##### ####        ##   ###       ####       ######     ####       ###  ###       ####       ######
      ###      #####    #####  ###    #####  ####       ###  ####  ###    ##### #####    #####  ###       #### #####  #####     ####       ###  #####   ###### ####   #####
      ###       ############   ###     ##### ####       ###  ####  ###     ##### #############  ###       #### ############     ####       ###  ############## ############
      ###          ######      ###       ###  ###       ###  ####  ###       ###    ######      ###       ###     ######         ###       ###     ###### ###     ######                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      
	`)

	//ical.GetCalendarEvents("https://outlook.office365.com/owa/calendar/941374dceabe4f82a2fe81bb7437b459@teknikenshus.se/040d9692705a435fb6c32078aa5a1130911275803346469963/calendar.ics")
	// Connect to the MQTT broker
	f, err := os.Open("calendar.ics")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	// Calculate the first and last day of the current month
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	start := firstDay
	end := lastDay
	fmt.Println("Start: ", start)
	fmt.Println("End: ", end)

	calendar := gocal.NewParser(f)
	calendar.Start, calendar.End = &start, &end
	err = calendar.Parse()
	if err != nil {
		log.Fatalf("Error parsing calendar: %v", err)
	}

	fmt.Println("Events for the month:")
	for _, event := range calendar.Events {

		fmt.Printf("%s on %s Location: %s", event.Summary, event.Start, event.Location)
		fmt.Println()
	}
	fmt.Println()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1)
	fmt.Println("Today Start: ", todayStart)
	fmt.Println("Today End: ", todayEnd)

	var todayEvents []gocal.Event
	for _, event := range calendar.Events {
		if event.Start.After(todayStart) && event.Start.Before(todayEnd) {
			todayEvents = append(todayEvents, event)
		}
	}

	fmt.Println("Events for today:")
	for _, event := range todayEvents {
		fmt.Printf("%s on %s Location: %s", event.Summary, event.Start, event.Location)
		fmt.Println()
	}
	fmt.Println()

	mqttClient, err := mqtt.NewClient()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}
	mqtt.Publish(mqttClient, "testTopic", "Hello World!", false)
	for i, event := range todayEvents {
		summaryTopic := fmt.Sprintf("sirius/today/%02d/summary", i+1)
		startTopic := fmt.Sprintf("sirius/today/%02d/start", i+1)
		endTopic := fmt.Sprintf("sirius/today/%02d/end", i+1)
		mqtt.Publish(mqttClient, summaryTopic, event.Summary, false)
		mqtt.Publish(mqttClient, startTopic, event.Start.Format("15:04"), false)
		mqtt.Publish(mqttClient, endTopic, event.End.Format("15:04"), false)

		fmt.Println()
	}
}
