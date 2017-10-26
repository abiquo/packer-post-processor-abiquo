package abiquo_api

type LocationCollection struct {
	AbstractCollection
	Collection []Location
}

type Location struct {
	DTO
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
