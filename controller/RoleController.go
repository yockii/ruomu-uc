package controller

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-uc/domain"
	"gorm.io/gorm"
	"strconv"

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

func (c *roleController) GetRoleResourceCodes(_ map[string][]string, value []byte) (any, error) {
	roleId := gjson.GetBytes(value, "roleId").Int()
	if roleId == 0 {
		return nil, nil
	}
	// 获取用户对应的权限和角色
	var resources []*model.Resource
	subSql := database.DB.Model(&model.RoleResource{}).Select("resource_id").Where("role_id=?", roleId)
	err := database.DB.Select("id, resource_code").Where("id in (?)", subSql).Find(&resources).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	var resourceCodes []string
	for _, resource := range resources {
		resourceCodes = append(resourceCodes, resource.ResourceCode)
	}
	return &shared.AuthorizationInfo{
		ResourceCodes: resourceCodes,
	}, nil
}

func (_ *roleController) Add(_ map[string][]string, value []byte) (any, error) {
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

	var c int64
	if err := database.DB.Model(&model.Role{}).Where(&model.Role{RoleName: instance.RoleName}).Count(&c).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	if c > 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		}, nil
	}

	instance.ID = util.SnowflakeId()
	if err := database.DB.Create(instance).Error; err != nil {
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

func (_ *roleController) Update(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Model(&model.Role{ID: instance.ID}).Updates(&model.Role{
		RoleName: instance.RoleName,
		RoleDesc: instance.RoleDesc,
	}).Error; err != nil {
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

func (_ *roleController) Delete(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Where(instance).Delete(instance).Error; err != nil {
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

func (_ *roleController) Instance(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.Role)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}
	if err := database.DB.Where(instance).First(instance).Error; err != nil {
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

func (_ *roleController) DispatchResourcesToRole(_ map[string][]string, value []byte) (any, error) {
	instance := new(domain.RoleResources)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.RoleID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " roleId",
		}, nil
	}
	if len(instance.ResourceIDs) == 0 {
		// 删除所有
		if err := database.DB.Where(&model.RoleResource{RoleID: instance.RoleID}).Delete(&model.RoleResource{}).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	} else {
		var oldResourceIds []uint64
		if err := database.DB.Model(&model.RoleResource{}).Where("role_id=?", instance.RoleID).Pluck("resource_id", &oldResourceIds).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
		// 要删除和添加的
		var toDelete []uint64
		var toAdd []uint64
		for _, oldId := range oldResourceIds {
			if !containsUint64(instance.ResourceIDs, oldId) {
				toDelete = append(toDelete, oldId)
			}
		}
		for _, newId := range instance.ResourceIDs {
			if !containsUint64(oldResourceIds, newId) {
				toAdd = append(toAdd, newId)
			}
		}
		// 执行数据库操作
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			if len(toDelete) > 0 {
				if err := tx.Where("role_id=? and resource_id in (?)", instance.RoleID, toDelete).Delete(&model.RoleResource{}).Error; err != nil {
					return err
				}
			}
			if len(toAdd) > 0 {
				var roleResources []model.RoleResource
				for _, id := range toAdd {
					roleResources = append(roleResources, model.RoleResource{
						ID:         util.SnowflakeId(),
						RoleID:     instance.RoleID,
						ResourceID: id,
					})
				}
				if err := tx.Create(&roleResources).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	}

	// 清除缓存
	conn := cache.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Errorln(err)
		}
	}(conn)
	_, err := conn.Do("DEL", shared.RedisKeyRoleResourceCode+strconv.FormatUint(instance.RoleID, 10))
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		}, nil
	}

	return &server.CommonResponse{Data: true}, nil
}

func (_ *roleController) List(_ map[string][]string, value []byte) (any, error) {
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

	tx := database.DB.Limit(paginate.Limit).Offset(paginate.Offset)

	condition := &model.Role{
		ID: instance.ID,
	}
	if instance.RoleName != "" {
		tx.Where("role_name like ?", "%"+instance.RoleName+"%")
		instance.RoleName = ""
	}
	var total int64
	var list []*model.Role
	err := tx.Find(&list, condition).Offset(-1).Count(&total).Error
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
