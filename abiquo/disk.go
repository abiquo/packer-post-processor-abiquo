package abiquo_api

type DiskCollection struct {
	AbstractCollection
	Collection []Disk
}

type Disk struct {
	DTO
	Label              string `json:"label,omitempty"`
	Sequence           int    `json:"sequence,omitempty"`
	Path               string `json:"path,omitempty"`
	DiskFormatType     string `json:"diskFormatType,omitempty"`
	DiskFileSize       int    `json:"diskFileSize,omitempty"`
	HdRequired         int    `json:"hdRequired,omitempty"`
	State              string `json:"state,omitempty"`
	DiskControllerType string `json:"diskControllerType,omitempty"`
	DiskController     string `json:"diskController,omitempty"`
	CreationDate       string `json:"creationDate,omitempty"`
	Bootable           bool   `json:"bootable,omitempty"`
}
