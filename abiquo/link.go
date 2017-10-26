package abiquo_api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty"
)

type Link struct {
	Type  string `json:"type,omitempty"`
	Href  string `json:"href,omitempty"`
	Title string `json:"title,omitempty"`
	Rel   string `json:"rel,omitempty"`
}

func (l *Link) Get(c *AbiquoClient) (*resty.Response, error) {
	resp, err := c.checkResponse(c.client.R().SetHeader("Accept", l.Type).Get(l.Href))
	return resp, err
}

func (l *Link) trimPort() {
	r, _ := url.Parse(l.Href)

	var trimport bool
	if r.Scheme == "https" && r.Port() == "443" {
		trimport = true
	} else if r.Scheme == "http" && r.Port() == "80" {
		trimport = true
	} else {
		trimport = false
	}
	if trimport {
		l.Href = fmt.Sprintf("%s://%s%s?%s", r.Scheme, r.Hostname(), r.Path, r.RawQuery)
		l.Href = strings.Trim(l.Href, "?")
	}
}
