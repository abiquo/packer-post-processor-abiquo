package abiquo_api

type HWprofile struct {
	DTO
	Name   string `json:"name,omitempty"`
	Cpu    int    `json:"cpu,omitempty"`
	Ram    int    `json:"ramInMb,omitempty"`
	Active bool   `json:"active,omitempty"`
}

type HWprofileCollection struct {
	AbstractCollection
	Collection []HWprofile
}
