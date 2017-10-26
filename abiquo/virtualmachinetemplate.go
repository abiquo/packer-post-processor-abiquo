package abiquo_api

import (
	"encoding/json"
	"fmt"
)

type TemplateCollection struct {
	AbstractCollection
	Collection []VirtualMachineTemplate
}

type VirtualMachineTemplate struct {
	DTO
	Name                             string `json:"name,omitempty"`
	ChefEnabled                      bool   `json:"chefEnabled,omitempty"`
	CpuRequired                      int    `json:"cpuRequired,omitempty"`
	CreationDate                     string `json:"creationDate,omitempty"`
	CreationUser                     string `json:"creationUser,omitempty"`
	Description                      string `json:"description,omitempty"`
	EthernetDriverType               string `json:"ethernetDriverType,omitempty"`
	IconUrl                          string `json:"iconUrl,omitempty"`
	Id                               int    `json:"id,omitempty"`
	LoginPassword                    string `json:"loginPassword,omitempty"`
	LoginUser                        string `json:"loginUser,omitempty"`
	OsType                           string `json:"osType,omitempty"`
	OsVersion                        string `json:"osVersion,omitempty"`
	RamRequired                      int    `json:"ramRequired,omitempty"`
	State                            string `json:"state,omitempty"`
	EnableCpuHotAdd                  bool   `json:"enableCpuHotAdd,omitempty"`
	EnableRamHotAdd                  bool   `json:"enableRamHotAdd,omitempty"`
	EnableDisksHotReconfigure        bool   `json:"enableDisksHotReconfigure,omitempty"`
	EnableNicsHotReconfigure         bool   `json:"enableNicsHotReconfigure,omitempty"`
	EnableRemoteAccessHotReconfigure bool   `json:"enableRemoteAccessHotReconfigure,omitempty"`
}

func (t *VirtualMachineTemplate) GetDisks(c *AbiquoClient) ([]Disk, error) {
	var disksCol DiskCollection
	var disks []Disk

	disks_resp, err := t.FollowLink("disks", c)
	if err != nil {
		return disks, err
	}
	json.Unmarshal(disks_resp.Body(), &disksCol)

	for {
		for _, d := range disksCol.Collection {
			disks = append(disks, d)
		}
		if disksCol.HasNext() {
			next_link := disksCol.GetNext()
			disks_resp, err = c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.disks+json").
				Get(next_link.Href))
			if err != nil {
				return disks, err
			}
			json.Unmarshal(disks_resp.Body(), &disksCol)
		} else {
			break
		}
	}

	return disks, nil
}

func (t *VirtualMachineTemplate) ReplacePrimaryDisk(c *AbiquoClient, diskdef DiskDef, file string) (VirtualMachineTemplate, error) {
	var newTemplate VirtualMachineTemplate

	diskdef_json, err := json.Marshal(diskdef)
	if err != nil {
		return newTemplate, err
	}

	template_lnk, _ := t.GetLink("templatePath")
	templateUpdateUrl := template_lnk.Href

	params := map[string]string{
		"diskInfo": string(diskdef_json),
	}
	resp, err := c.upload(templateUpdateUrl, params, "diskFile", file)
	if err != nil {
		return newTemplate, err
	}
	if resp.StatusCode > 399 {
		err = fmt.Errorf("ERROR %s - HTTP %d", resp.Status, resp.StatusCode)
		return newTemplate, err
	}

	respref, err := t.Refresh(c)
	if err != nil {
		return newTemplate, err
	}
	json.Unmarshal(respref.Body(), &newTemplate)
	return newTemplate, err
}

func (t *VirtualMachineTemplate) Update(c *AbiquoClient) error {
	edit_lnk, _ := t.GetLink("edit")
	resp, err := c.checkResponse(c.client.R().SetHeader("Accept", edit_lnk.Type).
		SetHeader("Content-Type", edit_lnk.Type).
		Put(edit_lnk.Href))

	if err != nil {
		return err
	}
	json.Unmarshal(resp.Body(), t)
	return nil
}
