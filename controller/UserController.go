package controller

import (
	"encoding/json"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"

	"github.com/yockii/ruomu-uc/model"
	"github.com/yockii/ruomu-uc/service"
)

var UserController = new(userController)

type userController struct{}

func (c *userController) GetUserRoleIds(value []byte) (any, error) {
	uid := gjson.GetBytes(value, "userId").Int()
	if uid == 0 {
		return nil, nil
	}
	// 获取用户对应的权限和角色
	var roles []*model.Role
	err := database.DB.Where("id in (select role_id from t_user_role where user_id=?)", uid).Find(&roles)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	var roleIds []string

	for _, role := range roles {
		if role.RoleType == 99 {
			roleIds = append(roleIds, shared.SuperAdmin)
		} else {
			roleIds = append(roleIds, strconv.FormatInt(role.Id, 10))
		}
	}
	return &shared.AuthorizationInfo{
		RoleIds: roleIds,
	}, nil
}

func (c *userController) Add(value []byte) (interface{}, error) {
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

func (c *userController) Login(value []byte) (any, error) {
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
	if has, err := database.DB.Get(user); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	} else if !has {
		return &server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		}, nil
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(instance.Password)); err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch,
		}, nil
	}

	jwtToken, err := generateJwtToken(strconv.FormatInt(user.Id, 10), "")
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

func (c *userController) Update(value []byte) (any, error) {
	instance := new(model.User)
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

	if _, err := database.DB.ID(instance.Id).Update(&model.User{
		RealName:     instance.RealName,
		ExternalId:   instance.ExternalId,
		ExternalType: instance.ExternalType,
		Status:       instance.Status,
	}); err != nil {
		return &server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		}, nil
	}

	return &server.CommonResponse{Data: true}, nil
}

func (c *userController) Delete(value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if _, err := database.DB.Delete(instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	return &server.CommonResponse{Data: true}, nil
}

func (c *userController) Instance(value []byte) (any, error) {
	instance := new(model.User)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if has, err := database.DB.Get(instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	} else if !has {
		return &server.CommonResponse{}, nil
	}
	return &server.CommonResponse{Data: instance}, nil
}

func (c *userController) List(value []byte) (any, error) {
	instance := new(model.User)
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

	condition := &model.User{
		Id:           instance.Id,
		ExternalId:   instance.ExternalId,
		ExternalType: instance.ExternalType,
		Status:       instance.Status,
	}
	if instance.Username != "" {
		session.Where("username like ?", "%"+instance.Username+"%")
		instance.Username = ""
	}
	if instance.RealName != "" {
		session.Where("real_name like ?", "%"+instance.RealName+"%")
		instance.RealName = ""
	}

	var list []*model.User
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

func generateJwtToken(userId, tenantId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sid := util.GenerateXid()

	conn := cache.Get()
	defer conn.Close()
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

//
//func (_ *userController) Paginate(ctx *fiber.Ctx) error {
//	type UserCondition struct {
//		model.User
//		CreateTimeRange *server.TimeCondition `json:"createTimeRange"`
//	}
//	pr := new(UserCondition)
//	if err := ctx.QueryParser(pr); err != nil {
//		logger.Error(err)
//		return ctx.JSON(&server.CommonResponse{
//			Code: server.ResponseCodeParamParseError,
//			Msg:  server.ResponseMsgParamParseError,
//		})
//	}
//	limit, offset, orderBy, err := server.ParsePaginationInfoFromQuery(ctx)
//	if err != nil {
//		logger.Error(err)
//		return ctx.JSON(&server.CommonResponse{
//			Code: server.ResponseCodeParamParseError,
//			Msg:  server.ResponseMsgParamParseError,
//		})
//	}
//
//	timeRangeMap := make(map[string]*server.TimeCondition)
//	if pr.CreateTimeRange != nil {
//		timeRangeMap["update_time"] = &server.TimeCondition{
//			Start: pr.CreateTimeRange.Start,
//			End:   pr.CreateTimeRange.End,
//		}
//	}
//
//	total, list, err0 := service.UserService.PaginateBetweenTimes(&pr.User, limit, offset, orderBy, timeRangeMap)
//	if err0 != nil {
//		logger.Error(err0)
//		return ctx.JSON(&server.CommonResponse{
//			Code: server.ResponseCodeDatabase,
//			Msg:  server.ResponseMsgDatabase,
//		})
//	}
//	return ctx.JSON(&server.CommonResponse{Data: &server.Paginate{
//		Total:  total,
//		Offset: offset,
//		Limit:  limit,
//		Items:  list,
//	}})
//}
