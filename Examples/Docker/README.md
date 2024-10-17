## Using Docker Compose
Use the provided example `docker-compose.yaml` file as a starting point:

```yaml
services:
  calendar-to-mqtt:
    image: ghcr.io/teknikens-hus/calendar-to-mqtt:latest
    environment:
      TZ: "Europe/Stockholm"
    volumes:
      - ./config.yaml:/config.yaml
    restart: unless-stopped
```
Adjust the values in the config.yaml.example file and then rename it to config.yaml

To see options for the config file, check the main [README.md](../../README.md)

Then run:
```bash
docker-compose up -d
```
## Using Docker command:
```bash
docker run -d \
  --name calendar-to-mqtt \
  -e TZ="Europe/Stockholm" \
  -v ./config.yaml:/config.yaml \
  --restart unless-stopped \
  ghcr.io/teknikens-hus/calendar-to-mqtt:latest
```