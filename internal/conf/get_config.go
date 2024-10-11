package conf

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type MQTTConfig struct {
	ServerAddress string
	ClientID      string
	Username      string
	Password      string
	QoS           int
	WriteToLog    bool
}

type ConfigFile struct {
	MQTT struct {
		BrokerIP string `yaml:"BrokerIP"`
		ClientID string `yaml:"ClientID"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		QoS      int    `yaml:"QoS"`
		Log      bool   `yaml:"Log"`
	} `yaml:"MQTT"`
}

func GetMQTTConfig() *MQTTConfig {
	fmt.Println("Getting MQTT Config")
	config, err := readConf()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	fmt.Println("Creating MQTT Config")
	if config.MQTT.BrokerIP == "" {
		log.Fatalf("BrokerIP is not set in the config file")
	}
	// Create MQTT Config from the config file
	return &MQTTConfig{
		ServerAddress: config.MQTT.BrokerIP,
		ClientID:      config.MQTT.ClientID,
		Username:      config.MQTT.Username,
		Password:      config.MQTT.Password,
		QoS:           config.MQTT.QoS,
		WriteToLog:    config.MQTT.Log,
	}
}

func readConf() (*ConfigFile, error) {
	// The config file should be created in the root of the project
	fmt.Println("Reading config file")
	const filename = "config.yaml"
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	configPointer := &ConfigFile{}
	err = yaml.Unmarshal(buf, configPointer)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return configPointer, nil
}
