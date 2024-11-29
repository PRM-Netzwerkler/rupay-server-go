package util

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	mqttClient MQTT.Client
}

var (
	instance *Client
	once     sync.Once
)

func NewClient(con Config) error {
	// This method creates some default options for us, most notably it sets the auto reconnect option to be true, and the default port to `1883`. Auto reconnect is really useful in IOT applications as the internet connection may not always be extremely strong.
	opts := MQTT.NewClientOptions()

	// The specified The connection type we are using is just plain unencrypted TCP/IP
	opts.AddBroker(con.MqttBroker)
	// The client id needs to be unique, The argument passed was generated through a random number generator to avoid collisions.
	opts.SetClientID("con.MqttClientId")

	opts.SetUsername(con.MqttClientName)
	opts.SetPassword(con.MqttClientPassword)

	tlsConfig, _ := createTLSConfig(con)

	mqttClient := MQTT.NewClient(opts)

	opts.SetTLSConfig(tlsConfig)

	// Set automatic reconnect
	opts.SetAutoReconnect(true)

	// We have to create the connection to the broker manually and verify that there is no error.
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	print("Connected to mqtt client")

	instance = &Client{
		mqttClient,
	}

	return nil
}

func createTLSConfig(con Config) (*tls.Config, error) {

	// Decode the Base64 encoded cert and key
	caCert, err := base64.StdEncoding.DecodeString(con.CertCaRoot)
	if err != nil {
		log.Fatalf("Failed to decode CA certificate: %v", err)
	}

	mosquittoCert, err := base64.StdEncoding.DecodeString(con.CertMosquitto)
	if err != nil {
		log.Fatalf("Failed to decode Mosquitto certificate: %v", err)
	}

	mosquittoKey, err := base64.StdEncoding.DecodeString(con.KeyMosquitto)
	if err != nil {
		log.Fatalf("Failed to decode Mosquitto key: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// load client cert and key
	cert, err := tls.X509KeyPair(mosquittoCert, mosquittoKey)
	if err != nil {
		panic("Failed to load client certificate")
	}

	tlsConfig := &tls.Config{RootCAs: caCertPool, Certificates: []tls.Certificate{cert}, InsecureSkipVerify: false}

	return tlsConfig, nil
}

// GetClient returns the global MQTT client instance.
// Ensure `InitializeClient` is called before using this function.
func GetClient() *Client {
	if instance == nil {
		panic("MQTT client is not initialized. Call InitializeClient first.")
	}
	return instance
}

// Publish publishes a message on a specific topic. An error is returned if there was problem. This function will publish with a QOS of 1.
func (c *Client) Publish(msg, topic string) error {
	if token := c.mqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Subscribe creates a subscription for the passed topic. The callBack function is used to process any messages that the client recieves on that topic. The subscription created will have a QOS of 1.
func (c *Client) Subsribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// WaitForMessage subscribes to a topic and waits for a single message with an optional timeout.
// It returns the received message or an error.
func (c *Client) WaitForMessage(topic string, timeout time.Duration) (string, error) {
	messageChan := make(chan string)
	errorChan := make(chan error)

	// Define a handler to capture the message
	handler := func(client MQTT.Client, msg MQTT.Message) {
		messageChan <- string(msg.Payload())
	}

	// Subscribe to the topic
	if token := c.mqttClient.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		return "", token.Error()
	}
	defer c.mqttClient.Unsubscribe(topic) // Unsubscribe after receiving the message

	select {
	case message := <-messageChan:
		return message, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout waiting for message on topic: %s", topic)
	case err := <-errorChan:
		return "", err
	}
}
