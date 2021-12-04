package tibber

import (
	"github.com/gorilla/websocket"
	"github.com/tskaard/tibber-golang"
	"log"
	"time"
)

type Tibber struct {
	accessToken       string
	homeId            string
	wss               *websocket.Conn
	LatestMeasurement LiveMeasurement
}

func NewTibber(accessToken string, homeId string) (*Tibber, error) {
	return &Tibber{
		accessToken: accessToken,
		homeId:      homeId,
	}, nil
}

func (t *Tibber) Measurements() {

	s := tibber.NewStream(t.homeId, t.accessToken)
	msg := make(tibber.MsgChan)
	err := s.StartSubscription(msg)
	if err != nil {
		log.Fatalln(err)
	}

	func(msgChan tibber.MsgChan) {
		for {
			select {
			case msg2 := <-msgChan:

				t.LatestMeasurement.Timestamp = msg2.Payload.Data.LiveMeasurement.Timestamp
				t.LatestMeasurement.Power = msg2.Payload.Data.LiveMeasurement.Power
				t.LatestMeasurement.LastMeterConsumption = msg2.Payload.Data.LiveMeasurement.LastMeterConsumption
				t.LatestMeasurement.LastMeterProduction = msg2.Payload.Data.LiveMeasurement.LastMeterProduction
				t.LatestMeasurement.AccumulatedConsumption = msg2.Payload.Data.LiveMeasurement.AccumulatedConsumption
				t.LatestMeasurement.AccumulatedProduction = msg2.Payload.Data.LiveMeasurement.AccumulatedProduction
				t.LatestMeasurement.AccumulatedConsumptionLastHour = 0
				t.LatestMeasurement.AccumulatedProductionLastHour = 0
				t.LatestMeasurement.AccumulatedCost = msg2.Payload.Data.LiveMeasurement.AccumulatedCost
				t.LatestMeasurement.AccumulatedReward = msg2.Payload.Data.LiveMeasurement.AccumulatedReward
				t.LatestMeasurement.Currency = 0
				t.LatestMeasurement.MinPower = msg2.Payload.Data.LiveMeasurement.MinPower
				t.LatestMeasurement.AveragePower = msg2.Payload.Data.LiveMeasurement.AveragePower
				t.LatestMeasurement.MaxPower = msg2.Payload.Data.LiveMeasurement.MaxPower
				t.LatestMeasurement.PowerProduction = msg2.Payload.Data.LiveMeasurement.PowerProduction
				t.LatestMeasurement.PowerReactive = 0
				t.LatestMeasurement.PowerProductionReactive = 0
				t.LatestMeasurement.MinPowerProduction = msg2.Payload.Data.LiveMeasurement.MinPowerProduction
				t.LatestMeasurement.MaxPowerProduction = msg2.Payload.Data.LiveMeasurement.MaxPowerProduction
				t.LatestMeasurement.PowerFactor = 0

				updateFloatIfNotNull(&t.LatestMeasurement.VoltagePhase1, msg2.Payload.Data.LiveMeasurement.VoltagePhase1)
				updateFloatIfNotNull(&t.LatestMeasurement.VoltagePhase2, msg2.Payload.Data.LiveMeasurement.VoltagePhase2)
				updateFloatIfNotNull(&t.LatestMeasurement.VoltagePhase3, msg2.Payload.Data.LiveMeasurement.VoltagePhase3)

				updateFloatIfNotNull(&t.LatestMeasurement.CurrentL1, msg2.Payload.Data.LiveMeasurement.CurrentPhase1)
				updateFloatIfNotNull(&t.LatestMeasurement.CurrentL2, msg2.Payload.Data.LiveMeasurement.CurrentPhase2)
				updateFloatIfNotNull(&t.LatestMeasurement.CurrentL3, msg2.Payload.Data.LiveMeasurement.CurrentPhase3)
				t.LatestMeasurement.SignalStrength = 0
			}
		}
	}(msg)
}

func updateFloatIfNotNull(old *float64, new float64) {
	if new != 0 {
		*old = new
	}
}

type Msg struct {
	Type    string  `json:"type"`
	HomeID  string  `json:"homeId"`
	ID      int     `json:"id"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Data Data `json:"data"`
}

type Data struct {
	LiveMeasurement `json:"liveMeasurement"`
}

type LiveMeasurement struct {
	Timestamp                      time.Time `json:"timestamp"`
	Power                          float64   `json:"power"`
	LastMeterConsumption           float64   `json:"lastMeterConsumption"`
	LastMeterProduction            float64   `json:"lastMeterProduction"`
	AccumulatedConsumption         float64   `json:"accumulatedConsumption"`
	AccumulatedProduction          float64   `json:"accumulatedProduction"`
	AccumulatedConsumptionLastHour float64   `json:"accumulatedConsumptionLastHour"`
	AccumulatedProductionLastHour  float64   `json:"accumulatedProductionLastHour"`
	AccumulatedCost                float64   `json:"accumulatedCost"`
	AccumulatedReward              float64   `json:"accumulatedReward"`
	Currency                       float64   `json:"currency"`
	MinPower                       float64   `json:"minPower"`
	AveragePower                   float64   `json:"averagePower"`
	MaxPower                       float64   `json:"maxPower"`
	PowerProduction                float64   `json:"powerProduction"`
	PowerReactive                  float64   `json:"powerReactive"`
	PowerProductionReactive        float64   `json:"powerProductionReactive"`
	MinPowerProduction             float64   `json:"minPowerProduction"`
	MaxPowerProduction             float64   `json:"maxPowerProduction"`
	PowerFactor                    float64   `json:"powerFactor"`
	VoltagePhase1                  float64   `json:"voltagePhase1"`
	VoltagePhase2                  float64   `json:"voltagePhase2"`
	VoltagePhase3                  float64   `json:"voltagePhase3"`
	CurrentL1                      float64   `json:"currentL1"`
	CurrentL2                      float64   `json:"currentL2"`
	CurrentL3                      float64   `json:"currentL3"`
	SignalStrength                 float64   `json:"signalStrength"`
}

func (t *Tibber) FetchLiveMeasurement() LiveMeasurement {
	return t.LatestMeasurement
}
