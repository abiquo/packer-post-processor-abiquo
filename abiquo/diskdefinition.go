package abiquo_api

type DiskDef struct {
	DTO
	Bootable                  bool   `json:"bootable,omitempty"`
	Sequence                  int    `json:"sequence,omitempty"`
	RequiredHDInMB            int    `json:"requiredHDInMB,omitempty"`
	DiskFileFormat            string `json:"diskFileFormat,omitempty"`
	VirtualMachineTemplateUrl string `json:"virtualMachineTemplateUrl,omitempty"`
	DiskUrl                   string `json:"diskUrl,omitempty"`
	CurrentPath               string `json:"currentPath,omitempty"`
}
