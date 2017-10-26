package abiquo_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type VMTDef struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	CategoryName       string `json:"categoryName,omitempty"`
	DiskFileFormat     string `json:"diskFileFormat,omitempty"`
	RequiredCpu        string `json:"requiredCpu,omitempty"`
	RequiredHDInMB     string `json:"requiredHDInMB,omitempty"`
	RequiredRamInMB    string `json:"requiredRamInMB,omitempty"`
	LoginUser          string `json:"loginUser,omitempty"`
	LoginPassword      string `json:"loginPassword,omitempty"`
	OsType             string `json:"osType,omitempty"`
	OsVersion          string `json:"osVersion,omitempty"`
	EthernetDriverType string `json:"ethernetDriverType,omitempty"`
}

func (def *VMTDef) Upload(c *AbiquoClient, repo Repo, file string) (VirtualMachineTemplate, error) {
	var newTemplate VirtualMachineTemplate

	definition_json, err := json.Marshal(def)
	if err != nil {
		return newTemplate, err
	}

	post_url := ""
	repo_uri_lnk, err := repo.GetLink("applianceManagerRepositoryUri")
	if err != nil {
		return newTemplate, err
	}
	post_url = repo_uri_lnk.Href + "/templates"

	params := map[string]string{
		"diskInfo": string(definition_json),
	}
	resp, err := c.upload(post_url, params, "diskFile", file)
	if err != nil {
		return newTemplate, err
	}

	location, err := resp.Location()
	if err != nil {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("Upload response did not had a location header! Upload did not succeed. Response body was '%s'", bodyBytes)
		return newTemplate, errors.New(errorMsg)
	}

	respt, err := c.checkResponse(c.client.R().
		SetHeader("Accept", "application/vnd.abiquo.virtualmachinetemplate+json").
		Get(location.String()))
	if err != nil {
		return newTemplate, err
	}

	json.Unmarshal(respt.Body(), &newTemplate)
	return newTemplate, err
}
