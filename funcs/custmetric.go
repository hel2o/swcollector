package funcs

import (
	"errors"
	"log"
	"strconv"

	"time"

	go_snmp "github.com/gosnmp/gosnmp"
	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type CustM struct {
	Ip           string
	custmMetrics []CustmMetric
}
type CustmMetric struct {
	metric     string
	tag        string
	value      float64
	metrictype string
}

func AllCustomIp(ipRange []string) (allIp []string) {
	if len(ipRange) > 0 {
		for _, sip := range ipRange {
			aip := sw.ParseIp(sip)
			for _, ip := range aip {
				allIp = append(allIp, ip)
			}
		}
	}
	return allIp
}

func CustomMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	if !g.Config().CustomMetrics.Enabled {
		return
	}
	chs := make([]chan CustM, 0)
	for _, ip := range AliveIp {
		if ip != "" {
			for _, metric := range g.CustConfig().Metrics {
				customIps := AllCustomIp(metric.IpRange)
				if g.InArray(ip, customIps) {
					ch := make(chan CustM)
					go custMetrics(ip, metric, ch)
					chs = append(chs, ch)
				}
			}

		}
	}
	for _, ch := range chs {
		custom, ok := <-ch
		if !ok {
			continue
		}

		for _, customMetric := range custom.custmMetrics {
			L = append(L, GaugeValueIp(time.Now().Unix(), custom.Ip, SwcollectorTakeSec, time.Since(startTime).Seconds(), "type=custom"))
			if customMetric.metrictype == "GAUGE" {
				L = append(L, GaugeValueIp(time.Now().Unix(), custom.Ip, customMetric.metric, customMetric.value, customMetric.tag))
			}
			if customMetric.metrictype == "COUNTER" {
				L = append(L, CounterValueIp(time.Now().Unix(), custom.Ip, customMetric.metric, customMetric.value, customMetric.tag))
			}
		}

	}
	endTime := time.Now()

	log.Printf("Update Custmetric complete. Process time %s.", endTime.Sub(startTime))

	return L
}

func custMetrics(ip string, metric *g.MetricConfig, ch chan CustM) {
	var custm CustM
	var custmmetric CustmMetric
	var custmmetrics []CustmMetric
	startTime := time.Now()

	value, err := GetCustMetric(ip, metric.Oid, g.Config().Switch.SnmpTimeout, g.Config().Switch.SnmpRetry)
	endTime := time.Now()
	if g.Config().Debug {
		log.Printf("ip: %s. metric: %s. Process time %s.", ip, metric.Metric, endTime.Sub(startTime))
	}
	if err != nil {
		log.Println(ip, metric.Oid, err)
		close(ch)
		return
	} else {
		custmmetric.metric = metric.Metric
		custmmetric.metrictype = metric.Type
		custmmetric.tag = metric.Tag
		custmmetric.value = value
		custmmetrics = append(custmmetrics, custmmetric)
	}

	custm.Ip = ip
	custm.custmMetrics = custmmetrics
	ch <- custm
	return
}

func GetCustMetric(ip, oid string, timeout, retry int) (float64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in CustomMetric, Oid is ", oid, r)
		}
	}()
	method := "get"
	var value float64
	var err error
	var snmpPDUs []go_snmp.SnmpPDU
	snmpPDUs, err = sw.RunSnmp(ip, g.GetCommunity(ip), oid, method, retry, timeout)
	if len(snmpPDUs) > 0 {
		value, err = interfaceTofloat64(snmpPDUs[0].Value)
	}
	return value, err
}

func interfaceTofloat64(v interface{}) (float64, error) {
	var err error
	switch value := v.(type) {
	case int:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case uint8:
		return float64(value), nil
	case uint16:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		value_parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		} else {
			return value_parsed, nil
		}
	default:
		err = errors.New("value cannot not Parse to digital")
		return 0, err
	}
}
