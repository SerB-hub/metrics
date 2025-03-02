package storage

type (
	Gauge   float64
	Counter int64
)

type Storage interface {
	SaveGaugeMetricValue(
		name string,
		value Gauge,
	) error
	SaveCounterMetricValue(
		name string,
		value Counter,
	) error
	GetGaugeMetricValue(name string) (Gauge, error)
	GetCounterMetricValue(name string) (Counter, error)
}
