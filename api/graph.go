package api

import (
	"log"

	"gopkg.in/macaron.v1"
)

// GetGraph will return the link of grafana pannel
func (c *Client) GetGraph(ctx *macaron.Context) string {
	appID := ctx.Query("appid")
	metrics := ctx.Query("metrics")
	from := ctx.Query("from")
	to := ctx.Query("to")

	width := func(w string) string {
		if w == "" {
			w = "900"
		}
		return w
	}(ctx.Query("width"))

	height := func(h string) string {
		if h == "" {
			h = "400"
		}
		return h
	}(ctx.Query("height"))

	link, err := c.generateLink(appID, metrics, from, to, width, height)
	if err != nil {
		return "Faild to get the link of Grafana."
	}
	return link
}

func (c *Client) generateLink(appID, metrics, from, to, width, height string) (string, error) {
	const (
		unixTime = "unixnano"
	)

	grafanaAddr := c.GrafServer
	dashbord := "dataman"
	panelID := c.setPanelID(metrics)

	start, end, err := timeRange(from, to, unixTime)
	if err != nil {
		log.Println(err)
		return "The from time is after the to time.", err
	}

	link := "<iframe src=\"" + grafanaAddr + "/dashboard-solo/db/" + dashbord + "?panelId=" + panelID + "orTab=Axes&from=" + start + "&to=" + end + "&var-AppID=" + appID + "\" width=\"" + width + "\" height=\"" + height + "\" frameborder=\"0\"></iframe>\n"
	log.Printf("[Info] Get the grafana panel link: %s", link)
	return link, nil
}

func (c *Client) setPanelID(metrics string) (panelID string) {
	switch metrics {
	case "cpuusage":
		panelID = "1"
	case "memusage":
		panelID = "3"
	default:
		panelID = "1"
	}
	return panelID
}
