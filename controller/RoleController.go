package controller

import (
	"encoding/json"

	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"

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

func (_ *roleController) Add(value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.RoleName == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name",
		}, nil
	}

	if c, err := database.DB.Count(&model.Role{RoleName: instance.RoleName}); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	} else if c > 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		}, nil
	}

	instance.Id = util.SnowflakeId()
	if _, err := database.DB.Insert(instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: instance,
	}, nil
}

func (_ *roleController) Update(value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.Id == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if _, err := database.DB.Update(&model.Role{
		RoleName: instance.RoleName,
		RoleDesc: instance.RoleDesc,
	}, &model.Role{Id: instance.Id}); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *roleController) Delete(value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.Id == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if _, err := database.DB.Delete(instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *roleController) List(value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	paginate := new(server.Paginate)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	session := database.DB.NewSession().Limit(paginate.Limit, paginate.Offset)

	condition := &model.Role{
		Id: instance.Id,
	}
	if instance.RoleName != "" {
		session.Where("role_name like ?", "%"+instance.RoleName+"%")
		instance.RoleName = ""
	}

	var list []*model.Role
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Offset: paginate.Offset,
			Limit:  paginate.Limit,
			Items:  list,
		},
	}, nil
}
