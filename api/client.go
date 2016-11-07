package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/promproxy/util"
)

const (
	queryPath      = "/api/v1/query"
	queryRangePath = "/api/v1/query_range"
	timeFormat     = "2006-01-02 15:04:05"
)

// Client is a client for excuting queries against the Prometheus API and Grafana link
type Client struct {
	PromServer string
	GrafServer string
	httpClient *http.Client
}

// NewClient creats a new Client, given the PromServer, GrafServer
func NewClient(conf util.Conf) *Client {
	return &Client{
		PromServer: conf.PromServer,
		GrafServer: conf.GrafServer,
		httpClient: http.DefaultClient,
	}
}

func timeRange(f, t, unixTime string) (string, string, error) {
	if f == "" && t == "" {
		to := time.Now()
		from := timeOffset(to, "-30m")
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f == "" && t != "" {
		to, err := time.Parse(timeFormat, t)
		if err != nil {
			return "", "", err
		}
		from := timeOffset(to, "-30m")
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f != "" && t == "" {
		from, err := time.Parse(timeFormat, f)
		if err != nil {
			return "", "", err
		}
		to := timeOffset(from, "30m")
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	from, err := time.Parse(timeFormat, f)
	if err != nil {
		return "", "", err
	}
	to, err := time.Parse(timeFormat, t)
	if err != nil {
		return "", "", err
	}

	if to.Before(from) {
		return "", "", errors.New("The from time is after the to time.")
	}

	return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
}

func timeConvertString(t time.Time, unixTime string) string {
	var toUnix int64
	switch unixTime {
	case "unix":
		toUnix = t.Unix()
	case "unixnano":
		toUnix = t.UnixNano() / 1000000
	default:
		toUnix = t.Unix()
	}

	return strconv.FormatInt(toUnix, 10)
}

func timeOffset(t time.Time, offset string) time.Time {
	duration, _ := time.ParseDuration(offset)
	return t.Add(duration)
}
