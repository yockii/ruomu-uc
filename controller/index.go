package controller

import (
	"encoding/json"

	"github.com/yockii/ruomu-core/shared"

	"github.com/yockii/ruomu-uc/constant"
)

func Dispatch(code string, headers map[string][]string, value []byte) ([]byte, error) {
	switch code {
	// 代码注入点
	case shared.InjectCodeAuthorizationInfoByUserId:
		return wrapCall(value, UserController.GetUserRoleIds)
	case shared.InjectCodeAuthorizationInfoByRoleId:
		return wrapCall(value, RoleController.GetRoleResourceCodes)
	// HTTP 注入点
	case constant.InjectCodeUserLogin:
		return wrapCall(value, UserController.Login)
	case constant.InjectCodeUserAdd:
		return wrapCall(value, UserController.Add)
	case constant.InjectCodeUserUpdate:
		return wrapCall(value, UserController.Update)
	case constant.InjectCodeUserDelete:
		return wrapCall(value, UserController.Delete)
	case constant.InjectCodeUserInstance:
		return wrapCall(value, UserController.Instance)
	case constant.InjectCodeUserList:
		return wrapCall(value, UserController.List)

	}
	return nil, nil
}

func wrapCall(v []byte, f func([]byte) (any, error)) ([]byte, error) {
	r, err := f(v)
	if err != nil {
		return nil, err
	}
	bs, err := json.Marshal(r)
	return bs, err
}
