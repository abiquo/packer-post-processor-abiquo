package abiquo_api

type ConfigProperty struct {
	DTO
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type ConfigPropertyCollection struct {
	AbstractCollection
	Collection []ConfigProperty
}
