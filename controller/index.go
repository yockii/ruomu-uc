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
		return wrapCall(headers, value, UserController.GetUserRoleIds)
	case shared.InjectCodeAuthorizationInfoByRoleId:
		return wrapCall(headers, value, RoleController.GetRoleResourceCodes)
	// HTTP 注入点
	//// 用户
	case constant.InjectCodeUserLogin:
		return wrapCall(headers, value, UserController.Login)
	case constant.InjectCodeUserAdd:
		return wrapCall(headers, value, UserController.Add)
	case constant.InjectCodeUserUpdate:
		return wrapCall(headers, value, UserController.Update)
	case constant.InjectCodeUserDelete:
		return wrapCall(headers, value, UserController.Delete)
	case constant.InjectCodeUserInstance:
		return wrapCall(headers, value, UserController.Instance)
	case constant.InjectCodeUserList:
		return wrapCall(headers, value, UserController.List)
	case constant.InjectCodeUserPassword:
		return wrapCall(headers, value, UserController.UpdatePassword)
	//// 角色
	case constant.InjectCodeRoleAdd:
		return wrapCall(headers, value, RoleController.Add)
	case constant.InjectCodeRoleUpdate:
		return wrapCall(headers, value, RoleController.Update)
	case constant.InjectCodeRoleDelete:
		return wrapCall(headers, value, RoleController.Delete)
	case constant.InjectCodeRoleInstance:
		return wrapCall(headers, value, RoleController.Instance)
	case constant.InjectCodeRoleList:
		return wrapCall(headers, value, RoleController.List)

	}
	return nil, nil
}

func wrapCall(h map[string][]string, v []byte, f func(map[string][]string, []byte) (any, error)) ([]byte, error) {
	r, err := f(h, v)
	if err != nil {
		return nil, err
	}
	bs, err := json.Marshal(r)
	return bs, err
}
