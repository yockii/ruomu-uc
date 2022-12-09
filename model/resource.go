package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type Resource struct {
	Id           int64             `json:"id,omitempty" xorm:"pk"`
	ResourceName string            `json:"resourceName,omitempty" xorm:"comment('资源名称')"` // 资源名称
	ResourceCode string            `json:"resourceCode,omitempty" xorm:"comment('资源代码')"` // 资源认证代码
	HttpMethod   string            `json:"httpMethod,omitempty" xorm:"comment('http方法')"` // http方法
	CreateTime   database.DateTime `json:"createTime" xorm:"created"`
	UpdateTime   database.DateTime `json:"updateTime" xorm:"updated"`
}

func (_ Resource) TableComment() string {
	return "资源表"
}
func (r *Resource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.Id = j.Get("id").Int()
	r.ResourceName = j.Get("resourceName").String()
	r.ResourceCode = j.Get("resourceCode").String()
	r.HttpMethod = j.Get("httpMethod").String()
	return nil
}

type RoleResource struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	RoleId     int64             `json:"roleId,omitempty"`
	ResourceId int64             `json:"resourceId,omitempty"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (_ RoleResource) TableComment() string {
	return "角色资源表"
}
func (rr *RoleResource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	rr.Id = j.Get("id").Int()
	rr.RoleId = j.Get("roleId").Int()
	rr.ResourceId = j.Get("resourceId").Int()
	return nil
}
