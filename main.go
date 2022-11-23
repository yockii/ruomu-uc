package main

import (
	"encoding/json"

	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/shared"

	"github.com/yockii/ruomu-uc/controller"
)

type UC struct{}

func (UC) Initial(params map[string]string) error {
	for key, value := range params {
		config.Set(key, value)
	}
	return nil
}

func (UC) InjectCall(code string, value []byte) ([]byte, error) {
	switch code {
	case CodeUserAdd:
		user, err := controller.UserController.Add(value)
		if err != nil {
			return nil, err
		} else {
			bs, err := json.Marshal(user)
			return bs, err
		}
	}

	return nil, nil
}

func main() {
	shared.ModuleServe(ModuleName, &UC{})
}
