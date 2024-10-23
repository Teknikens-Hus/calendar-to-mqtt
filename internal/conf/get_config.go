package conf

import (
	"fmt"
	"os"
	"sync"

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
	ICS []ICSConfig `yaml:"ICS"`
}

type ICSConfig struct {
	Name     string `yaml:"Name"`
	URL      string `yaml:"URL"`
	Interval int    `yaml:"Interval"`
}

func GetMQTTConfig() *MQTTConfig {
	fmt.Println("Conf: Getting MQTT Config...")
	config, err := readConf()
	if err != nil {
		log.Fatalf("Conf: Error reading config file: %v", err)
	}
	if config.MQTT.BrokerIP == "" {
		log.Fatalf("Conf: BrokerIP is not set in the config file")
	}
	fmt.Println("Conf: Retrieved MQTT Configuration!")
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

func GetICSConfig() (*[]ICSConfig, error) {
	fmt.Println("Conf: Getting ICS Config...")
	config, err := readConf()
	if err != nil {
		log.Error("Conf: Error reading config file: %w", err)
		return nil, err
	}
	if len(config.ICS) == 0 {
		log.Info("Conf: No ICS configurations found in the config file")
		return nil, nil
	}
	fmt.Println("Conf: Retrieved ICS Configuration!")
	// Valid check all values for null
	for _, ics := range config.ICS {
		if ics.Name == "" {
			log.Fatalf("Conf: Name is missing in one ICS in the config file")
		}
		if ics.URL == "" {
			log.Fatalf("Conf: URL is missing in one ICS in the config file")
		}
		if ics.Interval == 0 || ics.Interval < 0 {
			log.Fatalf("Conf: Interval misisng in one ICS in the config file")
		}
	}
	// Create ICS Config from the config file
	return &config.ICS, nil
}

var (
	config     *ConfigFile
	configOnce sync.Once
)

func readConf() (*ConfigFile, error) {
	configOnce.Do(func() {
		// The config file should be created in the root of the project
		const filename = "config.yaml"
		fmt.Println("Conf: Reading config file: " + filename)
		buf, err := os.ReadFile(filename)
		if err != nil {
			return
		}

		configPointer := &ConfigFile{}
		err = yaml.Unmarshal(buf, configPointer)
		if err != nil {
			config = nil
			return
		}

		config = configPointer
	})

	if config == nil {
		return nil, fmt.Errorf("config file not found")
	}

	return config, nil
}
