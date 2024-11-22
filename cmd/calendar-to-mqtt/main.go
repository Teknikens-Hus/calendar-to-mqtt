// calendar-to-mqtt is a simple application that reads a Calendars (ICS etc) and publishes the events to an MQTT broker.
package main

import (
	"fmt"

	"os"
	"time"
	_ "time/tzdata"

	ics "github.com/Teknikens-Hus/calendar-to-mqtt/internal/calendars/ics"
	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/mqtt"
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

	// Manually update timezone from TZ env variable
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	fmt.Printf("Timezone set to: %s\n", time.Local)
	fmt.Printf("Current time: %s\n", time.Now().Format(time.RFC3339))

	// Connect to the MQTT broker

	mqttClient, err := mqtt.NewClient()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}

	fmt.Println()
	ics.SetupICS(&mqttClient)

	// Keep the application running
	select {}
}
