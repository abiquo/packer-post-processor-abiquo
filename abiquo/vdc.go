package abiquo_api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type VdcCollection struct {
	AbstractCollection
	Collection []VDC
}

type VDC struct {
	DTO
	HypervisorType    string `json:"hypervisorType,omitempty"`
	Name              string `json:"name,omitempty"`
	SyncState         string `json:"syncState,omitempty"`
	DiskSoftLimitInMb int    `json:"diskSoftLimitInMb,omitempty"`
	DiskHardLimitInMb int    `json:"diskHardLimitInMb,omitempty"`
	StorageSoftInMb   int    `json:"storageSoftInMb,omitempty"`
	StorageHardInMb   int    `json:"storageHardInMb,omitempty"`
	VlansSoft         int    `json:"vlansSoft,omitempty"`
	VlansHard         int    `json:"vlansHard,omitempty"`
	PublicIpsSoft     int    `json:"publicIpsSoft,omitempty"`
	PublicIpsHard     int    `json:"publicIpsHard,omitempty"`
	RAMSoft           int    `json:"ramSoft,omitempty"`
	RAMHard           int    `json:"ramHard,omitempty"`
	CPUSoft           int    `json:"cpuSoft,omitempty"`
	CPUHard           int    `json:"cpuHard,omitempty"`
}

func (v *VDC) GetVirtualApps(c *AbiquoClient) ([]VirtualApp, error) {
	var allVapps []VirtualApp
	var vapps VirtualAppCollection
	vapps_raw, err := v.FollowLink("virtualappliances", c)
	if err != nil {
		return allVapps, err
	}
	json.Unmarshal(vapps_raw.Body(), &vapps)
	for {
		for _, va := range vapps.Collection {
			allVapps = append(allVapps, va)
		}
		if vapps.HasNext() {
			next_link := vapps.GetNext()
			vapps_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualappliances+json").
				Get(next_link.Href))
			if err != nil {
				return allVapps, err
			}
			json.Unmarshal(vapps_raw.Body(), &vapps)
		} else {
			break
		}
	}
	return allVapps, nil
}

func (v *VDC) GetTemplate(template_name string, c *AbiquoClient) (VirtualMachineTemplate, error) {
	var vt VirtualMachineTemplate
	templates, err := v.GetTemplates(c)
	if err != nil {
		return vt, err
	}
	for _, t := range templates {
		if t.Name == template_name {
			return t, nil
		}
	}
	errorMsg := fmt.Sprintf("Template '%s' not found in VDC '%s'", template_name, v.Name)
	return vt, errors.New(errorMsg)
}

func (v *VDC) GetTemplates(c *AbiquoClient) ([]VirtualMachineTemplate, error) {
	var templates TemplateCollection
	var alltemplates []VirtualMachineTemplate

	templates_raw, err := v.FollowLink("templates", c)
	if err != nil {
		return alltemplates, err
	}

	json.Unmarshal(templates_raw.Body(), &templates)
	for {
		for _, t := range templates.Collection {
			alltemplates = append(alltemplates, t)
		}

		if templates.HasNext() {
			next_link := templates.GetNext()
			templates_raw, err = c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualmachinetemplates+json").
				Get(next_link.Href))
			if err != nil {
				return alltemplates, err
			}
			json.Unmarshal(templates_raw.Body(), &templates)
		} else {
			break
		}
	}

	return alltemplates, nil
}

func (v *VDC) GetHardwareProfiles(c *AbiquoClient) ([]HWprofile, error) {
	var allProfiles []HWprofile
	var hprofiles HWprofileCollection
	var location Location

	location_raw, err := v.FollowLink("location", c)
	if err != nil {
		return allProfiles, err
	}
	json.Unmarshal(location_raw.Body(), &location)

	profiles_raw, err := location.FollowLink("hardwareprofiles", c)
	if err != nil {
		return allProfiles, err
	}

	json.Unmarshal(profiles_raw.Body(), &hprofiles)
	for {
		for _, hp := range hprofiles.Collection {
			allProfiles = append(allProfiles, hp)
		}

		if hprofiles.HasNext() {
			next_link := hprofiles.GetNext()
			profiles_raw, err = c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.hardwareprofiles+json").
				Get(next_link.Href))
			if err != nil {
				return allProfiles, err
			}
			json.Unmarshal(profiles_raw.Body(), &hprofiles)
		} else {
			break
		}
	}

	return allProfiles, nil
}

func (v *VDC) CreateVapp(vapp_name string, c *AbiquoClient) (VirtualApp, error) {
	var vapp VirtualApp
	vapps_lnk, _ := v.GetLink("virtualappliances")

	vapp.Name = vapp_name
	jsonbytes, err := json.Marshal(vapp)
	if err != nil {
		return vapp, err
	}
	vapp_raw, err := c.checkResponse(c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualappliance+json").
		SetHeader("Content-Type", "application/vnd.abiquo.virtualappliance+json").
		SetBody(jsonbytes).
		Post(vapps_lnk.Href)))
	if err != nil {
		return vapp, err
	}
	json.Unmarshal(vapp_raw.Body(), &vapp)
	return vapp, nil
}
