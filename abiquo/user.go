package abiquo_api

import (
	"encoding/json"
)

type UserCollection struct {
	AbstractCollection
	Collectio []User
}

type User struct {
	DTO
	Nick         string `json:"nick,omitempty"`
	Name         string `json:"name,omitempty"`
	Surname      string `json:"surname,omitempty"`
	Description  string `json:"description,omitempty"`
	Email        string `json:"email,omitempty"`
	Locale       string `json:"locale,omitempty"`
	AuthType     string `json:"authType,omitempty"`
	Active       bool   `json:"active,omitempty"`
	PublicSSHKey string `json:"publicSshKey,omitempty"`
	FirstLogin   bool   `json:"firstLogin,omitempty"`
	Locked       bool   `json:"locked,omitempty"`
}

func (u *User) GetEnterprise(c *AbiquoClient) (Enterprise, error) {
	var ent Enterprise
	ent_resp, err := u.FollowLink("enterprise", c)
	if err != nil {
		return ent, err
	}
	json.Unmarshal(ent_resp.Body(), &ent)
	return ent, nil
}
