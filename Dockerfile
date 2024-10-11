FROM golang:1.23.2 AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/calendar-to-mqtt -o /calendar-to-mqtt

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /calendar-to-mqtt /calendar-to-mqtt


#USER nonroot:nonroot

ENTRYPOINT ["/calendar-to-mqtt"]