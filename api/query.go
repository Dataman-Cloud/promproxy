package api

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"gopkg.in/macaron.v1"
)

// Query will return the resulf by query from Prometheus
func (c *Client) Query(ctx *macaron.Context) []byte {
	appID := ctx.Query("appid")
	metrics := ctx.Query("metrics")
	timePrama := ctx.Query("time")

	result, err := c.queryFromProm(appID, metrics, timePrama)
	if err != nil {
		log.Fatalf("Failed to query the %s by error: %v", metrics, err)
		return nil
	}
	return result
}

func (c *Client) queryFromProm(appID, metrics, t string) ([]byte, error) {
	const (
		unixTime = "unix"
	)

	expr := c.setQueryExpr(metrics, appID)

	u, err := url.Parse(c.PromServer)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/") + queryPath
	q := u.Query()
	q.Set("query", expr)

	if t != "" {
		timestamp, err := time.Parse(timeFormat, t)
		if err != nil {
			return nil, err
		}
		q.Set("time", timeConvertString(timestamp, unixTime))
	}
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Info] Get the prometheus qurey result by url: %s", u.String())

	return result, nil
}

// setQueryExpr will return the expr of prometheus query
func (c *Client) setQueryExpr(metrics, appID string) (expr string) {
	switch metrics {
	case "cpuusage":
		expr = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='" + appID + "',id=~'/docker/.*', name=~'mesos.*'}[5m])) by (instance, name)"
	case "memusage":
		expr = "container_memory_usage_bytes{container_label_APP_ID='" + appID + "',id=~'/docker/.*', name=~'mesos.*'} / container_spec_memory_limit_bytes{container_label_APP_ID='nginx-stress', id=~'/docker/.*', name=~'mesos.*'}"
	default:
		expr = ""
	}
	return expr
}
