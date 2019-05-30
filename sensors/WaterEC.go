package sensors

import (
	"encoding/json"
	"log"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mqtt"
)

type waterECMessage struct {
	ElectricalConductivity float64 `json:"ecLevel"`
}

type WaterElectricalConductivity struct {
	electricalConducivity float64
}

const WaterECTopic = "/state/water_ec"

func (e *WaterElectricalConductivity) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(WaterECTopic, e.ECHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *WaterElectricalConductivity) ECHandler(client mqttPaho.Client, message mqttPaho.Message) {
	eM := &waterECMessage{}
	json.Unmarshal(message.Payload(), eM)
	e.electricalConducivity = eM.ElectricalConductivity
}

func (e WaterElectricalConductivity) GetLatestData() float64 {
	return e.electricalConducivity
}