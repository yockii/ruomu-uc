package controller

import (
	"encoding/json"

	"github.com/yockii/ruomu-core/shared"

	"github.com/yockii/ruomu-uc/constant"
)

func Dispatch(code string, headers map[string]string, value []byte) ([]byte, error) {
	switch code {
	// 代码注入点
	case shared.InjectCodeAuthorizationInfoByUserId:
		return wrapCall(value, UserController.GetUserAuthorizationInfo)

	// HTTP 注入点
	case constant.CodeUserAdd:
		return wrapCall(value, UserController.Add)
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
