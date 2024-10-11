package mqtt

import (
	"fmt"
	"time"

	"github.com/Teknikens-Hus/calendar-to-mqtt/internal/conf"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	paho     paho.Client
	QoS      int
	ClientID string
}

func NewClient() (MQTTClient, error) {
	// paho mqtt is imported as paho instead of mqtt to distinguish it
	mqttconf := conf.GetMQTTConfig()
	opts := paho.NewClientOptions()
	opts.AddBroker(mqttconf.ServerAddress)
	opts.SetClientID(mqttconf.ClientID)

	opts.SetOrderMatters(false)       // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Set authentication if username and password is provided
	if mqttconf.Username != "" {
		fmt.Println("Setting username and password")
		opts.SetPassword(mqttconf.Password)
		opts.SetUsername(mqttconf.Username)
	}

	// Automate connection management (will keep trying to onnect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl paho.Client, err error) {
		fmt.Println("MQTT Connection lost")
	}
	opts.OnConnect = func(paho.Client) {
		fmt.Println("MQTT Connection established")
	}
	opts.OnReconnecting = func(paho.Client, *paho.ClientOptions) {
		fmt.Println("Attempting to reconnect to MQTT...")
	}

	//
	// Connect to the broker
	//
	client := paho.NewClient(opts)
	fmt.Println("Connecting to MQTT broker")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("MQTT Connection is up")

	// Return the client and the QoS level + error
	return MQTTClient{client, mqttconf.QoS, mqttconf.ClientID}, nil
}

func Publish(client MQTTClient, topic string, payload string, retain bool) {
	// Add clientID as prefix to topic
	topic = client.ClientID + "/" + topic
	token := client.paho.Publish(topic, byte(client.QoS), retain, payload)
	token.Wait()
	fmt.Printf("Published %s to topic: %s\n", payload, topic)
}
