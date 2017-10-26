package abiquo_api

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty"
)

// Generic DTO
type DTO struct {
	Links []Link `json:"links,omitempty"`
}

func (d *DTO) FollowLink(rel string, c *AbiquoClient) (*resty.Response, error) {
	link, err := d.GetLink(rel)
	if err != nil {
		return &resty.Response{}, err
	}

	resp, err := c.checkResponse(c.client.NewRequest().
		SetHeader("Accept", link.Type).
		Get(link.Href))
	return resp, err
}

func (d *DTO) GetLink(rel string) (Link, error) {
	link := Link{Href: ""}

	for _, l := range d.Links {
		if l.Rel == rel {
			link = l
		}
	}

	if link.Href == "" {
		errorMsg := fmt.Sprintf("Link with rel '%s' not found", rel)
		return link, errors.New(errorMsg)
	} else {
		link.trimPort()
		return link, nil
	}
}

func (d *DTO) Refresh(c *AbiquoClient) (*resty.Response, error) {
	edit_lnk, err := d.GetLink("edit")
	if err != nil {
		edit_lnk, _ = d.GetLink("self")
	}
	return c.checkResponse(c.client.R().SetHeader("Accept", edit_lnk.Type).
		SetHeader("Content-Type", edit_lnk.Type).
		Get(edit_lnk.Href))
}

// Generic Collection
type AbstractCollection struct {
	Links     []Link
	TotalSize int
}

func (c *AbstractCollection) GetNext() Link {
	link := Link{Href: ""}

	for _, l := range c.Links {
		if l.Rel == "next" {
			link = l
		}
	}

	if link.Href == "" {
		return Link{}
	} else {
		link.trimPort()
		return link
	}
}

func (c *AbstractCollection) HasNext() bool {
	for _, link := range c.Links {
		if link.Rel == "next" {
			return true
		}
	}
	return false
}
