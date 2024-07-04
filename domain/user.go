package domain

import "github.com/tidwall/gjson"

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserRoles struct {
	UserID  uint64   `json:"userId,string"`
	RoleIds []uint64 `json:"roleIds"`
}

func (ur *UserRoles) UnmarshalJSON(b []byte) error {
	ur.RoleIds = make([]uint64, 0)
	j := gjson.ParseBytes(b)
	ur.UserID = j.Get("userId").Uint()
	j.Get("roleIds").ForEach(func(_, v gjson.Result) bool {
		ur.RoleIds = append(ur.RoleIds, v.Uint())
		return true
	})
	return nil
}
