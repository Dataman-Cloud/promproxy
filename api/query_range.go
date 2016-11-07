package api

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"gopkg.in/macaron.v1"
)

//QueryRange will return the result by query_range from Prometheus
func (c *Client) QueryRange(ctx *macaron.Context) []byte {
	appID := ctx.Query("appid")
	metrics := ctx.Query("metrics")
	from := ctx.Query("from")
	to := ctx.Query("to")
	step := ctx.Query("step")

	result, err := c.queryRangeFromProm(appID, metrics, from, to, step)
	if err != nil {
		log.Fatalf("Failed to query range the %s by error: %v", metrics, err)
		return nil
	}
	return result
}

func (c *Client) queryRangeFromProm(appID, metrics, from, to, step string) ([]byte, error) {
	const (
		unixTime = "unix"
	)
	expr := c.setQueryExpr(metrics, appID)

	u, err := url.Parse(c.PromServer)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/") + queryRangePath
	q := u.Query()
	q.Set("query", expr)

	start, end, err := timeRange(from, to, unixTime)
	if err != nil {
		return nil, err
	}
	q.Set("start", start)
	q.Set("end", end)

	if step == "" {
		q.Set("step", "15s")
	} else {
		q.Set("step", step)
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
