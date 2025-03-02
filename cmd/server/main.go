package main

import (
	"fmt"
	"github.com/SerB-hub/metrics/internal/app/middlewares"
	storage "github.com/SerB-hub/metrics/internal/app/storage"
	"net/http"
	"strconv"
	"strings"
)

type (
	Gauge   float64
	Counter int64
)

type Config struct {
	MetricTypes map[string]string
}

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type UpdateMetricsHandler struct {
	config  *Config
	storage storage.Storage
}

func (umh *UpdateMetricsHandler) UpdateMetric(
	rw http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentTypeHeader := r.Header.Get("Content-Type"); !strings.HasPrefix(contentTypeHeader, "text/plain") {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	params := r.Context().Value("params").(map[string]string)
	metricType := params["metricType"]
	metricName := params["metricName"]
	metricValue := params["metricValue"]

	metricTypeExist := false
	metricTypeSame := false
	metricNameExist := false
	var saveError error

	for mtn, mt := range umh.config.MetricTypes {
		if strings.EqualFold(metricType, mtn) {
			metricTypeExist = true

			value, typeEqualCheck := umh.checkStringValueType(metricValue, mt)

			if !typeEqualCheck {
				break
			}

			metricTypeSame = true

			switch mtn {
			case storage.MetricTypeCounter:
				oldValue, err := umh.storage.GetCounterMetricValue(metricName)

				if err != nil {
					break
				}

				metricNameExist = true
				saveError = umh.storage.SaveCounterMetricValue(
					metricName,
					oldValue+storage.Counter((value).(int64)),
				)

				break

			case storage.MetricTypeGauge:
				_, err := umh.storage.GetGaugeMetricValue(metricName)

				if err != nil {
					break
				}

				metricNameExist = true
				saveError = umh.storage.SaveGaugeMetricValue(
					metricName,
					storage.Gauge((value).(float64)),
				)

				break
			}

			break
		}
	}

	if !metricTypeExist {
		fmt.Printf(`Metric type name "%s" does not exist\n`, metricType)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !metricTypeSame {
		fmt.Printf(`Value for metric type name "%s" is of the wrong type\n`, metricType)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !metricNameExist {
		fmt.Printf(`Metric name "%s" does not exist\n`, metricName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if saveError != nil {
		fmt.Printf(`%s\n`, saveError.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (umh *UpdateMetricsHandler) checkStringValueType(
	v string,
	t string,
) (interface{}, bool) {
	switch t {
	case "int64":
		if value, err := strconv.ParseInt(v, 10, 64); err == nil {
			return value, true
		}

		return nil, false

	case "float64":
		if value, err := strconv.ParseFloat(v, 64); err == nil {
			return value, true
		}

		return nil, false
	}

	return nil, false
}

func main() {
	config := &Config{
		MetricTypes: map[string]string{
			MetricTypeGauge:   "float64",
			MetricTypeCounter: "int64",
		},
	}
	memStorage := storage.NewMemStorage()
	memStorage.SaveGaugeMetricValue("g1", 1)
	updateMetricsHandler := &UpdateMetricsHandler{
		config:  config,
		storage: memStorage,
	}

	updateMetric := http.HandlerFunc(updateMetricsHandler.UpdateMetric)

	router := middlewares.NewRouter(
		map[string]*http.HandlerFunc{
			"/update/{metricType}/{metricName}/{metricValue}": &updateMetric,
		},
	)

	mux := http.NewServeMux()
	mux.Handle("/", router.ProcessRequest(nil))
	err := http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		panic(err)
	}
}
