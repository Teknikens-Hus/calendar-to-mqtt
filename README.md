# calendar-to-mqtt
This project aims to create a container that allows different types of calendars to be subscribed to and then published as MQTT topics

## NOTE!
This is currently not a working project. Under development!

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
```
If you dont use authentication to MQTT, just leave the username an empty string and it will be skipped.

## Supported Calendars
Currently, the following calendars are supported:
- None. TBC


## Development
1. Clone the repository
2. Install golang and or docker
3. Run `go mod download` to download the dependencies
4. Create a config.yaml file in the root of the project with the required values.
5. Run `go run /cmd/calendar-to-mqtt/main.go` to start the application or docker compose up --build to start the container
