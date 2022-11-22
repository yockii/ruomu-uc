package controller

import (
	"encoding/json"
	"errors"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/ruomu-core/util"

	"github.com/yockii/ruomu-uc/model"
	"github.com/yockii/ruomu-uc/service"
)

var UserController = new(userController)

type userController struct{}

func (_ *userController) Add(value []byte) (instance *model.User, err error) {
	instance = new(model.User)
	if err = json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return nil, err
	}

	// 处理必填
	if instance.Username == "" || (instance.Password == "" && instance.ExternalType == "") {
		err = errors.New("username and password or externalType is required")
		return
	}

	if instance.Password != "" {
		isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
		if !isStrong {
			return nil, errors.New("password strength is not enough")
		}
	}

	duplicated, success, err := service.UserService.Add(instance)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if duplicated {
		return nil, errors.New("user duplicated")
	}
	if success {
		return instance, nil
	}
	return nil, errors.New("Unknown error ")
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
