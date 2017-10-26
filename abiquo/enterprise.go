package abiquo_api

import (
	"encoding/json"
)

type EnterpriseCollection struct {
	AbstractCollection
	Collection []Enterprise
}

type Enterprise struct {
	DTO
	Name                             string `json:"name,omitempty"`
	IsReservationRestricted          bool   `json:"isReservationRestricted,omitempty"`
	Workflow                         bool   `json:"workflow,omitempty"`
	TwoFactorAuthenticationMandatory bool   `json:"twoFactorAuthenticationMandatory,omitempty"`
	DiskSoftLimitInMb                int    `json:"diskSoftLimitInMb,omitempty"`
	DiskHardLimitInMb                int    `json:"diskHardLimitInMb,omitempty"`
	StorageSoftInMb                  int    `json:"storageSoftInMb,omitempty"`
	StorageHardInMb                  int    `json:"storageHardInMb,omitempty"`
	VlansSoft                        int    `json:"vlansSoft,omitempty"`
	VlansHard                        int    `json:"vlansHard,omitempty"`
	PublicIpsSoft                    int    `json:"publicIpsSoft,omitempty"`
	PublicIpsHard                    int    `json:"publicIpsHard,omitempty"`
	RepositorySoftInMb               int    `json:"repositorySoftInMb,omitempty"`
	RepositoryHardInMb               int    `json:"repositoryHardInMb,omitempty"`
	RAMSoft                          int    `json:"ramSoft,omitempty"`
	RAMHard                          int    `json:"ramHard,omitempty"`
	CPUSoft                          int    `json:"cpuSoft,omitempty"`
	CPUHard                          int    `json:"cpuHard,omitempty"`
}

func (e *Enterprise) GetRepos(c *AbiquoClient) ([]Repo, error) {
	var reposcol RepoCollection
	var repos []Repo

	repos_resp, err := e.FollowLink("datacenterrepositories", c)
	if err != nil {
		return repos, err
	}
	json.Unmarshal(repos_resp.Body(), &reposcol)

	for {
		for _, r := range reposcol.Collection {
			repos = append(repos, r)
		}
		if reposcol.HasNext() {
			next_link := reposcol.GetNext()
			repos_raw, err := c.checkResponse(c.client.R().SetHeader("Accept", "application/vnd.abiquo.datacenterrepositories+json").
				Get(next_link.Href))
			if err != nil {
				return repos, err
			}
			json.Unmarshal(repos_raw.Body(), &reposcol)
		} else {
			break
		}
	}

	return repos, nil
}
