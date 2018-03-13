package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"os"
	"log"
	"io/ioutil"
)

const (
	location = "asia-east1"
	project  = "selfhydro-197504"
	registry = "raspberry-pis"
	device   = "original-hydro"
)

type MQTTComms struct {
	client MQTT.Client
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (mqtt *MQTTComms) authenticateDevice() {

	tokenString, _ := createJWTToken(project)

	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.googleapis.com:8883")

	clientId := "projects/" + project + "/locations/" + location + "/registries/" + registry + "/devices/" + device

	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(f)
	opts.SetPassword(tokenString)
	opts.SetProtocolVersion(4)
	opts.SetUsername("unused")

	mqtt.client = MQTT.NewClient(opts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
func (mqtt *MQTTComms) subscribeToTopic(topic string) {
	if token := mqtt.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
func (mqtt *MQTTComms) unsubscribeFromTopic(topic string) {
	if token := mqtt.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	mqtt.client.Disconnect(250)
}
func (mqtt *MQTTComms) publishMessage(topic string, message string) {

	text := fmt.Sprintf("%v", message)
	token := mqtt.client.Publish(topic, 0, false, text)
	token.Wait()
}

func createJWTToken(projectId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"aud": projectId,
	})

	file, err := os.Open("/selfhydro/rsa_private.pem") // For read access.
	if err != nil {
		log.Fatal(err)
	}


	key, _ := ioutil.ReadFile(file.Name())

	rsaPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(key)


	tokenString, err := token.SignedString(rsaPrivateKey)

	fmt.Println(tokenString, err)
	return tokenString, err
}
