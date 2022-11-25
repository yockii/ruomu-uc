package main

import (
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"

	"github.com/yockii/ruomu-uc/constant"
	"github.com/yockii/ruomu-uc/controller"
	"github.com/yockii/ruomu-uc/model"
)

type UC struct{}

func (UC) Initial(params map[string]string) error {
	for key, value := range params {
		config.Set(key, value)
	}

	database.Initial()

	database.DB.Sync2(
		model.User{},
	)
	return nil
}

func (UC) InjectCall(code string, headers map[string]string, value []byte) ([]byte, error) {
	return controller.Dispatch(code, headers, value)
}

func main() {
	util.InitNode(1)
	defer database.Close()
	shared.ModuleServe(constant.ModuleName, &UC{})
}
