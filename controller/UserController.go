package controller

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/ruomu-uc/domain"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yockii/ruomu-uc/model"
	"github.com/yockii/ruomu-uc/service"
)

var UserController = new(userController)

type userController struct{}

func (c *userController) GetUserRoleIds(_ map[string][]string, value []byte) (any, error) {
	uid := gjson.GetBytes(value, "userId").Int()
	if uid == 0 {
		return nil, nil
	}
	// 获取用户对应的权限和角色
	var roles []*model.Role

	subSql := database.DB.Model(&model.UserRole{}).Select("role_id").Where("user_id=?", uid)

	err := database.DB.Where("id in (?)", subSql).Find(&roles).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	var roleIds []string

	for _, role := range roles {
		if role.RoleType == -1 {
			roleIds = append(roleIds, shared.SuperAdmin)
		} else {
			roleIds = append(roleIds, strconv.FormatUint(role.ID, 10))
		}
	}
	return &shared.AuthorizationInfo{
		RoleIds: roleIds,
	}, nil
}

func (c *userController) Add(_ map[string][]string, value []byte) (interface{}, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.Username == "" || (instance.Password == "" && instance.ExternalType == "") {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " username / password & externalType",
		}, nil
	}

	if instance.Password != "" {
		isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
		if !isStrong {
			return &server.CommonResponse{
				Code: server.ResponseCodePasswordStrengthInvalid,
				Msg:  server.ResponseMsgPasswordStrengthInvalid,
			}, nil
		}
	}

	duplicated, success, err := service.UserService.Add(instance)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	if duplicated {
		return &server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		}, nil
	}
	if success {
		return &server.CommonResponse{
			Data: instance,
		}, nil
	}
	return &server.CommonResponse{
		Code: server.ResponseCodeUnknownError,
		Msg:  server.ResponseMsgUnknownError,
	}, nil
}

func (c *userController) Login(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " username / password",
		}, nil
	}
	isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
	if !isStrong {
		return &server.CommonResponse{
			Code: server.ResponseCodePasswordStrengthInvalid,
			Msg:  server.ResponseMsgPasswordStrengthInvalid,
		}, nil
	}

	user := &model.User{Username: instance.Username}
	if err := database.DB.Take(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{
				Code: server.ResponseCodeDataNotExists,
				Msg:  server.ResponseMsgDataNotExists,
			}, nil
		} else {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(instance.Password)); err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch,
		}, nil
	}

	jwtToken, err := generateJwtToken(strconv.FormatUint(user.ID, 10), "")
	if err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration,
		}, nil
	}
	user.Password = ""
	return &server.CommonResponse{
		Data: map[string]interface{}{
			"token": jwtToken,
			"user":  user,
		},
	}, nil
}

func (c *userController) Update(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.User)
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

	if err := database.DB.Model(&model.User{ID: instance.ID}).Updates(&model.User{
		RealName:     instance.RealName,
		ExternalId:   instance.ExternalId,
		ExternalType: instance.ExternalType,
		Status:       instance.Status,
	}).Error; err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		}, nil
	}

	return &server.CommonResponse{Data: true}, nil
}

func (c *userController) Delete(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Where(instance).Delete(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{Data: true}, nil
}

func (c *userController) Instance(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Omit("password").Where(instance).Take(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{}, nil
		}
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	return &server.CommonResponse{Data: instance}, nil
}

func (c *userController) List(_ map[string][]string, value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	paginate := new(server.Paginate)
	if err := json.Unmarshal(value, paginate); err != nil {
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

	condition := &model.User{
		ID:           instance.ID,
		ExternalId:   instance.ExternalId,
		ExternalType: instance.ExternalType,
		Status:       instance.Status,
	}
	if instance.Username != "" {
		tx.Where("username like ?", "%"+instance.Username+"%")
		instance.Username = ""
	}
	if instance.RealName != "" {
		tx.Where("real_name like ?", "%"+instance.RealName+"%")
		instance.RealName = ""
	}

	var list []*model.User
	var total int64
	err := tx.Omit("password").Find(&list, condition).Offset(-1).Count(&total).Error
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

func (c *userController) UpdatePassword(header map[string][]string, value []byte) (any, error) {
	instance := new(domain.UpdateUserPasswordRequest)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	uidStr, has := header[shared.JwtClaimUserId]
	if !has || uidStr[0] == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " userId",
		}, nil
	}
	if instance.NewPassword == "" || instance.OldPassword == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " oldPassword / newPassword",
		}, nil
	}
	uid, _ := strconv.ParseUint(uidStr[0], 10, 64)
	if uid == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	userInstance := new(model.User)
	if err := database.DB.Model(&model.User{}).Where(&model.User{ID: uid}).First(userInstance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		}, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userInstance.Password), []byte(instance.OldPassword)); err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch,
		}, nil
	}
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(instance.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration,
		}, nil
	}
	if err = database.DB.Model(&model.User{ID: uid}).Update("password", string(encryptedPwd)).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		}, nil
	}

	return &server.CommonResponse{Data: true}, nil
}

func (c *userController) DispatchRolesToUser(_ map[string][]string, value []byte) (any, error) {
	// 先更新数据库，然后清除缓存
	instance := new(domain.UserRoles)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if instance.UserID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " userId",
		}, nil
	}

	if len(instance.RoleIds) == 0 {
		// 直接删除所有
		if err := database.DB.Where("user_id=?", instance.UserID).Delete(&model.UserRole{}).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase,
			}, nil
		}
	} else {

		// 获取用户原来的roleId列表
		var oldRoleIdList []uint64
		if err := database.DB.Model(&model.UserRole{}).Where("user_id=?", instance.UserID).Pluck("role_id", &oldRoleIdList).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase,
			}, nil
		}
		// 找到要删除的和要添加的roleId
		var toDelete []uint64
		var toAdd []uint64
		for _, roleId := range instance.RoleIds {
			if !containsUint64(oldRoleIdList, roleId) {
				toAdd = append(toAdd, roleId)
			}
		}
		for _, roleId := range oldRoleIdList {
			if !containsUint64(instance.RoleIds, roleId) {
				toDelete = append(toDelete, roleId)
			}
		}
		// 执行数据库操作
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for _, roleId := range toDelete {
				if err := tx.Where("user_id=? and role_id=?", instance.UserID, roleId).Delete(&model.UserRole{}).Error; err != nil {
					return err
				}
			}
			for _, roleId := range toAdd {
				if err := tx.Create(&model.UserRole{
					ID:     util.SnowflakeId(),
					UserID: instance.UserID,
					RoleID: roleId,
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase,
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
	// 删除用户ID对应的角色ID缓存
	_, err := conn.Do("DEL", shared.RedisKeyUserRoles+strconv.FormatUint(instance.UserID, 10))
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		}, nil
	}

	return &server.CommonResponse{Data: true}, nil
}

func generateJwtToken(userId, tenantId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sid := util.GenerateXid()

	conn := cache.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Errorln(err)
		}
	}(conn)
	_, err := conn.Do("SETEX", shared.RedisSessionIdKey+sid, config.GetInt("userTokenExpire"), userId)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	claims[shared.JwtClaimUserId] = userId
	claims[shared.JwtClaimTenantId] = tenantId
	claims[shared.JwtClaimSessionId] = sid

	t, err := token.SignedString([]byte(shared.JwtSecret))
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	return t, nil
}

func containsUint64(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
