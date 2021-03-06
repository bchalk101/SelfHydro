package sensors

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mocks"
	mqttMocks "github.com/selfhydro/selfhydro/mqtt/mocks"
	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldReturnCurrentTempAndHumidity(t *testing.T) {
	e := AmbientTemperature{
		temperature: 23.2,
	}
	currentTemperature := e.GetLatestData()
	assert.Equal(t, e.temperature, currentTemperature)
}

func Test_ShouldSubscribeToEnvironmentTopic(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockMQTTClient := &mqttMocks.MockMQTTClient{}
	mockMQTTMessage := &mocks.MockMQTTMessage{
		ReceivedPayload: []byte(`{"temperature":20.76101}`),
	}
	mockMQTT.On("SubscribeToTopic", string("/state/ambient_temperature"), mock.Anything).Run(func(args mock.Arguments) {
		args[1].(mqtt.MessageHandler)(mockMQTTClient, mockMQTTMessage)
	}).Return(nil)
	e := &AmbientTemperature{}
	e.Subscribe(mockMQTT)
	assert.Equal(t, e.temperature, 20.76101)
}

func Test_ShouldGetLatestSensorVoltage(t *testing.T) {
	e := AmbientTemperature{
		temperature: 12.2,
		Sensor:      Sensor{batteryVoltage: 3.4, id: 111},
	}
	assert.Equal(t, e.GetLatestBatteryVoltage(), 3.4)
}
