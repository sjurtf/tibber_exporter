package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sjurtf/tibber_exporter/tibber"
	"time"
)

type tibberCollector struct {
	client *tibber.Tibber

	// https://developer.tibber.com/docs/reference
	// LiveMeasurement
	// ignoring timestamp, 26 metrics in total
	power                          *prometheus.Desc
	lastMeterConsumption           *prometheus.Desc
	accumulatedConsumption         *prometheus.Desc
	accumulatedProduction          *prometheus.Desc
	accumulatedConsumptionLastHour *prometheus.Desc
	accumulatedProductionLastHour  *prometheus.Desc
	accumulatedCost                *prometheus.Desc
	accumulatedReward              *prometheus.Desc
	currency                       *prometheus.Desc
	minPower                       *prometheus.Desc
	averagePower                   *prometheus.Desc
	maxPower                       *prometheus.Desc
	powerProduction                *prometheus.Desc
	powerReactive                  *prometheus.Desc
	powerProductionReactive        *prometheus.Desc
	minPowerProduction             *prometheus.Desc
	maxPowerProduction             *prometheus.Desc
	lastMeterProduction            *prometheus.Desc
	powerFactor                    *prometheus.Desc
	voltagePhase1                  *prometheus.Desc
	voltagePhase2                  *prometheus.Desc
	voltagePhase3                  *prometheus.Desc
	currentL1                      *prometheus.Desc
	currentL2                      *prometheus.Desc
	currentL3                      *prometheus.Desc
	signalStrength                 *prometheus.Desc
}

func NewTibberCollector(tibber *tibber.Tibber) prometheus.Collector {
	return &tibberCollector{
		client: tibber,
		power: prometheus.NewDesc(
			"tibber_power",
			"Consumption at the moment (Watt)",
			nil, nil),
		lastMeterConsumption: prometheus.NewDesc(
			"tibber_last_meter_consumption",
			"Last meter active import register state (kWh)",
			nil, nil),
		accumulatedConsumption: prometheus.NewDesc(
			"tibber_accumulated_consumption",
			"kWh consumed since midnight",
			nil, nil),
		accumulatedProduction: prometheus.NewDesc(
			"tibber_accumulated_production",
			"net kWh produced since midnight",
			nil, nil),
		accumulatedConsumptionLastHour: prometheus.NewDesc(
			"tibber_accumulated_consumption_last_hour",
			"kWh consumed since since last hour shift",
			nil, nil),
		accumulatedProductionLastHour: prometheus.NewDesc(
			"tibber_accumulated_production_last_hour",
			"net kWh produced since last hour shift",
			nil, nil),
		accumulatedCost: prometheus.NewDesc(
			"tibber_accumulated_cost",
			"Accumulated cost since midnight; requires active Tibber power deal",
			nil, nil),
		accumulatedReward: prometheus.NewDesc(
			"tibber_accumulated_reward",
			"Accumulated reward since midnight; requires active Tibber power deal",
			nil, nil),
		currency: prometheus.NewDesc(
			"tibber_currency",
			"Currency of displayed cost; requires active Tibber power deal",
			nil, nil),
		minPower: prometheus.NewDesc(
			"tibber_min_power",
			"Min consumption since midnight (Watt)",
			nil, nil),
		averagePower: prometheus.NewDesc(
			"tibber_average_power",
			"Average consumption since midnight (Watt)",
			nil, nil),
		maxPower: prometheus.NewDesc(
			"tibber_max_power",
			"Peak consumption since midnight (Watt)",
			nil, nil),
		powerProduction: prometheus.NewDesc(
			"tibber_power_production",
			"Net production (A-) at the moment (Watt)",
			nil, nil),
		powerReactive: prometheus.NewDesc(
			"tibber_power_reactive",
			"Reactive consumption (Q+) at the moment (kVAr)",
			nil, nil),
		powerProductionReactive: prometheus.NewDesc(
			"tibber_power_production_reactive",
			"Net reactive production (Q-) at the moment (kVAr)",
			nil, nil),
		minPowerProduction: prometheus.NewDesc(
			"tibber_min_power_production",
			"Min net production since midnight (Watt)",
			nil, nil),
		maxPowerProduction: prometheus.NewDesc(
			"tibber_max_power_production",
			"Max net production since midnight (Watt)",
			nil, nil),
		lastMeterProduction: prometheus.NewDesc(
			"tibber_last_meter_production",
			"Last meter active export register state (kWh)",
			nil, nil),
		powerFactor: prometheus.NewDesc(
			"tibber_power_factor",
			"Power factor (active power / apparent power)",
			nil, nil),
		voltagePhase1: prometheus.NewDesc(
			"tibber_voltage_phase_1",
			"Voltage on phase 1",
			nil, nil),
		voltagePhase2: prometheus.NewDesc(
			"tibber_voltage_phase_2",
			"Voltage on phase 2",
			nil, nil),
		voltagePhase3: prometheus.NewDesc(
			"tibber_voltage_phase_3",
			"Voltage on phase 3",
			nil, nil),
		currentL1: prometheus.NewDesc(
			"tibber_current_line_1",
			"Current on L1",
			nil, nil),
		currentL2: prometheus.NewDesc(
			"tibber_current_line_2",
			"Current on L2",
			nil, nil),
		currentL3: prometheus.NewDesc(
			"tibber_current_line_3",
			"Current on L3",
			nil, nil),
		signalStrength: prometheus.NewDesc(
			"tibber_signal_strength",
			"Device signal strength (Pulse - dB; Watty - percent)",
			nil, nil),
	}
}

func (t *tibberCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (t *tibberCollector) Collect(ch chan<- prometheus.Metric) {
	l := t.client.FetchLiveMeasurement()

	// Data not ready yet
	if l.Timestamp.IsZero() {
		return
	}

	// Data is old
	now := time.Now()
	if l.Timestamp.Before(now.Add(-time.Minute * 5)) {
		return
	}

	ch <- prometheus.MustNewConstMetric(t.power, prometheus.GaugeValue, l.Power)
	ch <- prometheus.MustNewConstMetric(t.lastMeterConsumption, prometheus.GaugeValue, l.LastMeterConsumption)
	ch <- prometheus.MustNewConstMetric(t.lastMeterProduction, prometheus.GaugeValue, l.LastMeterProduction)
	ch <- prometheus.MustNewConstMetric(t.accumulatedConsumption, prometheus.GaugeValue, l.AccumulatedConsumption)
	ch <- prometheus.MustNewConstMetric(t.accumulatedProduction, prometheus.GaugeValue, l.AccumulatedProduction)
	ch <- prometheus.MustNewConstMetric(t.accumulatedConsumptionLastHour, prometheus.GaugeValue, l.AccumulatedConsumptionLastHour)
	ch <- prometheus.MustNewConstMetric(t.accumulatedProductionLastHour, prometheus.GaugeValue, l.AccumulatedProductionLastHour)
	ch <- prometheus.MustNewConstMetric(t.accumulatedCost, prometheus.GaugeValue, l.AccumulatedCost)
	ch <- prometheus.MustNewConstMetric(t.accumulatedReward, prometheus.GaugeValue, l.AccumulatedReward)
	ch <- prometheus.MustNewConstMetric(t.currency, prometheus.GaugeValue, l.Currency)
	ch <- prometheus.MustNewConstMetric(t.minPower, prometheus.GaugeValue, l.MinPower)
	ch <- prometheus.MustNewConstMetric(t.averagePower, prometheus.GaugeValue, l.AveragePower)
	ch <- prometheus.MustNewConstMetric(t.maxPower, prometheus.GaugeValue, l.MaxPower)
	ch <- prometheus.MustNewConstMetric(t.powerProduction, prometheus.GaugeValue, l.PowerProduction)
	ch <- prometheus.MustNewConstMetric(t.powerReactive, prometheus.GaugeValue, l.PowerReactive)
	ch <- prometheus.MustNewConstMetric(t.powerProductionReactive, prometheus.GaugeValue, l.PowerProductionReactive)
	ch <- prometheus.MustNewConstMetric(t.minPowerProduction, prometheus.GaugeValue, l.MinPowerProduction)
	ch <- prometheus.MustNewConstMetric(t.maxPowerProduction, prometheus.GaugeValue, l.MaxPowerProduction)
	ch <- prometheus.MustNewConstMetric(t.powerFactor, prometheus.GaugeValue, l.PowerFactor)

	if l.VoltagePhase3 != 0 || l.VoltagePhase1 != 0 || l.VoltagePhase2 != 0 {
		ch <- prometheus.MustNewConstMetric(t.voltagePhase1, prometheus.GaugeValue, l.VoltagePhase1)
		ch <- prometheus.MustNewConstMetric(t.voltagePhase2, prometheus.GaugeValue, l.VoltagePhase2)
		ch <- prometheus.MustNewConstMetric(t.voltagePhase3, prometheus.GaugeValue, l.VoltagePhase3)
	}
	if l.CurrentL1 != 0 || l.CurrentL2 != 0 || l.CurrentL3 != 0 {
		ch <- prometheus.MustNewConstMetric(t.currentL1, prometheus.GaugeValue, l.CurrentL1)
		ch <- prometheus.MustNewConstMetric(t.currentL2, prometheus.GaugeValue, l.CurrentL2)
		ch <- prometheus.MustNewConstMetric(t.currentL3, prometheus.GaugeValue, l.CurrentL3)
	}
	ch <- prometheus.MustNewConstMetric(t.signalStrength, prometheus.GaugeValue, l.SignalStrength)

}
