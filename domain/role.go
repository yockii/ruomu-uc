package domain

import "github.com/tidwall/gjson"

type RoleResources struct {
	RoleID      uint64   `json:"roleId,omitempty,string"`
	ResourceIDs []uint64 `json:"resourceIds,omitempty"`
}

func (rr *RoleResources) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	rr.RoleID = j.Get("roleId").Uint()
	j.Get("resourceIds").ForEach(func(_, value gjson.Result) bool {
		rr.ResourceIDs = append(rr.ResourceIDs, value.Uint())
		return true
	})
	return nil
}
