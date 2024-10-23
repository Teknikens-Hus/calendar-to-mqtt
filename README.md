# calendar-to-mqtt
This program written in golang allows you to deploy a docker container/pod that scrapes/fetches data from supported calendars and publish them to MQTT topics for use in other systems, like ESPHome/Home Assistant.

## Supported Calendars
Currently, the following calendars are supported:
- ICS

## Installation
Currently amd64 and arm64 are supported.

### Docker / docker-compose:
[![Docker Icon](https://skillicons.dev/icons?i=docker&theme=light)](./Examples/Docker/README.md)

### Kubernetes Deployment:
[![Kube Icon](https://skillicons.dev/icons?i=kubernetes&theme=light)](./Examples/Kubernetes/README.md)

## Configuration
The configuration is done using a config.yaml file.
The file should be placed in the same directory as the docker-compose file (or in the root of the container)
The configuration file should have the following structure:
```yaml
MQTT:
  BrokerIP: "tcp://xxx.xxx.xx.xxx:1883"
  ClientID: "calendar-to-mqtt"
  Username: "username"
  Password: "password"
  QoS: 1
  Log: true

ICS:
  - Name: "Name Of Calendar 1"
    URL: "https://outlook.office365.com/owa/calendar/exammple@example.se/example/calendar1.ics"
    Interval: 60

  - Name: "Name Of Calendar 2"
    URL: "https://outlook.office365.com/owa/calendar/exammple@example.se/example/calendar2.ics"
    Interval: 60
```
### Check logs!
If you are having issues, check the logs of the application/container. It should give you some direction of whats wrong, which events are found and to what topics events are published to.

## Configuration Reference

| Key        | Description                                                                 | Example Value                          |
|------------|-----------------------------------------------------------------------------|----------------------------------------|
| **MQTT**   |                                                                             |                                        |
| BrokerIP   | The IP address and port of the MQTT broker.                                 | `tcp://xxx.xxx.xx.xxx:1883`            |
| ClientID   | The client ID to use when connecting to the MQTT broker.                    | `calendar-to-mqtt`                     |
| Username   | The username for authenticating with the MQTT broker.                       | `username`                             |
| Password   | The password for authenticating with the MQTT broker.                       | `password`                             |
| QoS        | The Quality of Service level for MQTT messages (0, 1, or 2).                | `1`                                    |
| Log        | Enable or disable logging.                                                  | `true`                                 |
| **ICS**    |                                                                             |                                        |
| Name       | The name of the calendar.                                                   | `Name Of Calendar 1`                   |
| URL        | The URL to the ICS file of the calendar.                                    | `https://outlook.office365.com/owa/calendar/exammple@example.se/example/calendar1.ics` |
| Interval   | The interval in seconds at which the calendar is fetched.                   | `60`                                   |

## MQTT Published topics
The topic format is:
- _ClientID/ICS - Name/today/all/events_
- _ClientID/ICS - Name/today/upcoming/events_

At the topic an array of json objects are published with the events for today.
#### Example:
Config.yaml
```yaml
MQTT:
  BrokerIP: "tcp://192.168.1.100:1883"
  ClientID: "great-success"
  Username: "Foo"
  Password: "Bar"
  QoS: 1
  Log: false

ICS:
  - Name: "Room 1"
    URL: "https://outlook.office365.com/xxxxx.ics"
    Interval: 60
```
The topics then becomes:
  - _great-success/Room 1/today/all/events_
  - _great-success/Room 1/today/upcoming/events_

All events contains all events for today.
Upcoming events contains all events minus the ones that have already passed.

A payload might look like:
```json
[
  {
    "summary":"Meeting with Dr. Emmett Brown",
    "start":"07:00",
    "end":"11:30"
  },
  {
    "summary":"Meeting with Marty McFly",
    "start":"14:00",
    "end":"16:10"
  }
]
```


## Development
1. Clone the repository
2. Install golang and or docker
3. Run `go mod download` to download the dependencies
4. Create a config.yaml file in the root of the project with the required values.
5. Run `go run /cmd/calendar-to-mqtt/main.go` to start the application or docker compose up --build to start the container
