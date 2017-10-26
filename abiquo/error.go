package abiquo_api

type Error struct {
	DTO
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorCollection struct {
	AbstractCollection
	Collection []Error
}
