package controller

import (
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"

	"github.com/yockii/ruomu-uc/model"
)

var RoleController = new(roleController)

type roleController struct{}

func (c *roleController) GetRoleResourceCodes(value []byte) (any, error) {
	roleId := gjson.GetBytes(value, "roleId").Int()
	if roleId == 0 {
		return nil, nil
	}
	// 获取用户对应的权限和角色
	var resources []*model.Resource
	err := database.DB.Cols("id").Where("id in (select resource_id from t_role_resource where role_id=?)", roleId).Find(&resources)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	var resourceCodes []string
	for _, role := range resources {
		resourceCodes = append(resourceCodes, role.ResourceCode)
	}
	return &shared.AuthorizationInfo{
		ResourceCodes: resourceCodes,
	}, nil
}
