package abiquo_api

type Task struct {
	DTO
	JobsExtended struct {
		DTO
		Collection []struct {
			Links         []interface{} `json:"links,omitempty"`
			ID            string        `json:"id,omitempty"`
			ParentTaskID  string        `json:"parentTaskId,omitempty"`
			Type          string        `json:"type,omitempty"`
			Description   string        `json:"description,omitempty"`
			State         string        `json:"state,omitempty"`
			RollbackState string        `json:"rollbackState,omitempty"`
			Timestamp     int           `json:"timestamp,omitempty"`
		} `json:"collection,omitempty"`
	} `json:"jobsExtended,omitempty"`
	TaskID    string `json:"taskId,omitempty"`
	UserID    string `json:"userId,omitempty"`
	Type      string `json:"type,omitempty"`
	OwnerID   string `json:"ownerId,omitempty"`
	State     string `json:"state,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}
