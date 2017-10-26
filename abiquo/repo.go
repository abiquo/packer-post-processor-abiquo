package abiquo_api

import (
	"encoding/json"
)

type RepoCollection struct {
	AbstractCollection
	Collection []Repo
}

type Repo struct {
	DTO
	Name               string
	RepositoryLocation string
}

func (r *Repo) GetTemplates(c *AbiquoClient) ([]VirtualMachineTemplate, error) {
	var templatesCol TemplateCollection
	var templates []VirtualMachineTemplate

	templates_resp, err := r.FollowLink("virtualmachinetemplates", c)
	if err != nil {
		return templates, err
	}
	json.Unmarshal(templates_resp.Body(), &templatesCol)

	for {
		for _, t := range templatesCol.Collection {
			templates = append(templates, t)
		}
		if templatesCol.HasNext() {
			next_link := templatesCol.GetNext()
			templates_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualmachinetemplates+json").
				SetQueryParam("master", "true").
				Get(next_link.Href))
			if err != nil {
				return templates, err
			}
			json.Unmarshal(templates_raw.Body(), &templatesCol)
		} else {
			break
		}
	}

	return templates, nil
}
