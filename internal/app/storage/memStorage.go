package storage

import (
	"errors"
	"fmt"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type MemStorage struct {
	storage map[string]map[string]interface{}
}

func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		storage: make(map[string]map[string]interface{}, 2),
	}

	ms.storage[MetricTypeGauge] = make(map[string]interface{})
	ms.storage[MetricTypeCounter] = make(map[string]interface{})

	return ms
}

func (ms *MemStorage) SaveGaugeMetricValue(
	name string,
	value Gauge,
) error {
	ms.storage[MetricTypeGauge][name] = value

	return nil
}

func (ms *MemStorage) SaveCounterMetricValue(
	name string,
	value Counter,
) error {
	ms.storage[MetricTypeCounter][name] = value

	return nil
}

func (ms *MemStorage) GetGaugeMetricValue(name string) (Gauge, error) {
	if value, ok := ms.storage[MetricTypeGauge][name]; ok {
		return value.(Gauge), nil
	}

	return 0, errors.New(fmt.Sprintf("Gauge metric with name %v not found", name))
}

func (ms *MemStorage) GetCounterMetricValue(name string) (Counter, error) {
	if value, ok := ms.storage[MetricTypeGauge][name]; ok {
		return value.(Counter), nil
	}

	return 0, errors.New(fmt.Sprintf("Counter metric with name %v not found", name))
}
