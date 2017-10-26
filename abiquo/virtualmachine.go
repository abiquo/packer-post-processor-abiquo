package abiquo_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type VirtualMachineCollection struct {
	AbstractCollection
	Collection []VirtualMachine
}

type VirtualMachine struct {
	DTO
	UUID              string                 `json:"uuid,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Label             string                 `json:"label,omitempty"`
	Description       string                 `json:"description,omitempty"`
	CPU               int                    `json:"cpu,omitempty"`
	RAM               int                    `json:"ram,omitempty"`
	VdrpEnabled       bool                   `json:"vdrpEnabled,omitempty"`
	VdrpPort          int                    `json:"vdrpPort,omitempty"`
	IDState           int                    `json:"idState,omitempty"`
	State             string                 `json:"state,omitempty"`
	IDType            int                    `json:"idType,omitempty"`
	Type              string                 `json:"type,omitempty"`
	HighDisponibility int                    `json:"highDisponibility,omitempty"`
	Password          string                 `json:"password,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Monitored         bool                   `json:"monitored,omitempty"`
	Protected         bool                   `json:"protected,omitempty"`
	Variables         map[string]string      `json:"variables,omitempty"`
	CreationTimestamp int64                  `json:"creationTimestamp,omitempty"`
	Backuppolicies    []interface{}          `json:"backuppolicies,omitempty"`
	LastSynchronize   int64                  `json:"lastSynchronize,omitempty"`
}

func (v *VirtualMachine) GetVapp(c *AbiquoClient) (VirtualApp, error) {
	var vapp VirtualApp
	vapp_raw, err := v.FollowLink("virtualappliance", c)
	if err != nil {
		return vapp, err
	}
	json.Unmarshal(vapp_raw.Body(), &vapp)
	return vapp, nil
}

func (v *VirtualMachine) Deploy(c *AbiquoClient) error {
	deploy_lnk, err := v.GetLink("deploy")
	accept_request_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.acceptedrequest+json").
		Post(deploy_lnk.Href))
	if err != nil {
		return err
	}
	var accept_request AcceptedRequest
	json.Unmarshal(accept_request_raw.Body(), &accept_request)

	for {
		vm_raw, err := v.Refresh(c)
		if err != nil {
			return err
		}
		json.Unmarshal(vm_raw.Body(), v)
		if v.State == "LOCKED" {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	task_lnk, _ := accept_request.GetLink("status")
	task_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.taskextended+json").
		Get(task_lnk.Href))
	if err != nil {
		return err
	}
	var task Task
	json.Unmarshal(task_raw.Body(), &task)
	if task.State != "FINISHED_SUCCESSFULLY" {
		errorMsg := fmt.Sprintf("Task to deploy VM %s failed. Check events.", v.Name)
		return errors.New(errorMsg)
	}
	return nil
}

func (v *VirtualMachine) PowerOn(c *AbiquoClient) error {
	return v.applyState("ON", c)
}

func (v *VirtualMachine) PowerOff(c *AbiquoClient) error {
	return v.applyState("OFF", c)
}

func (v *VirtualMachine) applyState(state string, c *AbiquoClient) error {
	body := fmt.Sprintf("{\"state\": \"%s\"}", state)
	state_lnk, _ := v.GetLink("state")
	accept_request_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.acceptedrequest+json").
		SetHeader("Content-Type", "application/vnd.abiquo.virtualmachinestate+json").
		SetBody(body).
		Put(state_lnk.Href))
	if err != nil {
		return err
	}
	var accept_request AcceptedRequest
	json.Unmarshal(accept_request_raw.Body(), &accept_request)

	for {
		vm_raw, err := v.Refresh(c)
		if err != nil {
			return err
		}
		json.Unmarshal(vm_raw.Body(), v)
		if v.State == "LOCKED" {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	task_lnk, _ := accept_request.GetLink("status")
	task_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.taskextended+json").
		Get(task_lnk.Href))
	if err != nil {
		return err
	}
	var task Task
	json.Unmarshal(task_raw.Body(), &task)
	if task.State != "FINISHED_SUCCESSFULLY" {
		errorMsg := fmt.Sprintf("Task to power %s VM %s failed. Check events.", state, v.Name)
		return errors.New(errorMsg)
	}
	return nil
}

func (v *VirtualMachine) Reset(c *AbiquoClient) error {
	body := ""
	reset_lnk, _ := v.GetLink("reset")
	accept_request_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.acceptedrequest+json").
		SetHeader("Content-Type", "application/vnd.abiquo.virtualmachinestate+json").
		SetBody(body).
		Post(reset_lnk.Href))
	if err != nil {
		return err
	}
	var accept_request AcceptedRequest
	json.Unmarshal(accept_request_raw.Body(), &accept_request)

	for {
		vm_raw, err := v.Refresh(c)
		if err != nil {
			return err
		}
		json.Unmarshal(vm_raw.Body(), v)
		if v.State == "LOCKED" {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	task_lnk, _ := accept_request.GetLink("status")
	task_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.taskextended+json").
		Get(task_lnk.Href))
	if err != nil {
		return err
	}
	var task Task
	json.Unmarshal(task_raw.Body(), &task)
	if task.State != "FINISHED_SUCCESSFULLY" {
		errorMsg := fmt.Sprintf("Task to reset VM %s failed. Check events.", v.Name)
		return errors.New(errorMsg)
	}
	return nil
}

func (v *VirtualMachine) Delete(c *AbiquoClient) error {
	edit_lnk, _ := v.GetLink("edit")
	_, err := c.checkResponse(c.client.R().Delete(edit_lnk.Href))
	if err != nil {
		return err
	}

	for {
		resp, err := v.Refresh(c)
		if err != nil {
			return err
		}
		if resp.StatusCode() == 404 {
			break
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}

func (v *VirtualMachine) GetIP() string {
	var nics []Link
	for _, l := range v.Links {
		if strings.HasPrefix(l.Rel, "nic") {
			nics = append(nics, l)
		}
	}

	// First external ips
	for _, n := range nics {
		if n.Type == "application/vnd.abiquo.externalip+json" {
			return n.Title
		}
	}
	// Then public
	for _, n := range nics {
		if n.Type == "application/vnd.abiquo.publicip+json" {
			return n.Title
		}
	}
	// And private
	for _, n := range nics {
		if n.Type == "application/vnd.abiquo.privateip+json" {
			return n.Title
		}
	}
	return ""
}

func (v *VirtualMachine) SetMetadata(mdata string, c *AbiquoClient) error {
	metadata_lnk, _ := v.GetLink("metadata")
	body, _ := json.Marshal(mdata)

	_, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.metadata+json").
		SetHeader("Content-Type", "application/vnd.abiquo.metadata+json").
		SetBody(body).
		Put(metadata_lnk.Href))
	if err != nil {
		return err
	}
	return nil
}

func (v *VirtualMachine) GetMetadata(c *AbiquoClient) (map[string]interface{}, error) {
	var mdata map[string]interface{}
	metadata, err := v.FollowLink("metadata", c)
	if err != nil {
		return mdata, err
	}

	json.Unmarshal(metadata.Body(), &mdata)
	return mdata, nil
}
