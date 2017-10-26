package abiquo_api

import (
	"encoding/json"
)

type VirtualAppCollection struct {
	AbstractCollection
	Collection []VirtualApp
}

type VirtualApp struct {
	DTO
	Error             int    `json:"error,omitempty"`
	HighDisponibility int    `json:"highDisponibility,omitempty"`
	Name              string `json:"name,omitempty"`
	PublicApp         int    `json:"publicApp,omitempty"`
	State             string `json:"state,omitempty"`
}

func (v *VirtualApp) Delete(c *AbiquoClient) error {
	edit_lnk, _ := v.GetLink("edit")
	_, err := c.checkResponse(c.client.R().Delete(edit_lnk.Href))
	if err != nil {
		return err
	}

	return nil

}

func (v *VirtualApp) GetVMs(c *AbiquoClient) ([]VirtualMachine, error) {
	var vms []VirtualMachine
	var vmsCol VirtualMachineCollection
	vms_raw, err := v.FollowLink("virtualmachines", c)
	if err != nil {
		return vms, err
	}
	json.Unmarshal(vms_raw.Body(), &vmsCol)

	for {
		for _, vm := range vmsCol.Collection {
			vms = append(vms, vm)
		}

		if vmsCol.HasNext() {
			next_link := vmsCol.GetNext()
			vms_raw, err = c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualmachines+json").
				Get(next_link.Href))
			if err != nil {
				return vms, err
			}
			json.Unmarshal(vms_raw.Body(), &vmsCol)
		} else {
			break
		}
	}
	return vms, nil
}

func (v *VirtualApp) CreateVM(vm VirtualMachine, c *AbiquoClient) (VirtualMachine, error) {
	var vm_created VirtualMachine
	body, _ := json.Marshal(vm)
	vms_lnk, _ := v.GetLink("virtualmachines")

	vm_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.virtualmachine+json").
		SetHeader("Content-Type", "application/vnd.abiquo.virtualmachine+json").
		SetBody(body).
		Post(vms_lnk.Href))
	if err != nil {
		return vm_created, err
	}
	json.Unmarshal(vm_raw.Body(), &vm_created)
	return vm_created, nil
}
