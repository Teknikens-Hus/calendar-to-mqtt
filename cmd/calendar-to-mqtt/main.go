// calendar-to-mqtt is a simple application that reads a Calendars (ICS etc) and publishes the events to an MQTT broker.
package main

import (
	"fmt"

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

	// Connect to the MQTT broker

	mqttClient, err := mqtt.NewClient()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}
	mqtt.Publish(mqttClient, "testTopic", "Hello World!", false)

	ics.SetupICS(&mqttClient)

	// Keep the application running
	select {}
}
