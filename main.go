package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"

	"github.com/yockii/ruomu-uc/constant"
	"github.com/yockii/ruomu-uc/controller"
	"github.com/yockii/ruomu-uc/model"
	"github.com/yockii/ruomu-uc/service"
)

type UC struct{}

func (UC) Initial(params map[string]string) error {
	for key, value := range params {
		config.Set(key, value)
	}

	database.Initial()

	database.DB.Sync2(
		model.User{},
		model.Role{},
		model.UserExtend{},
		model.UserRole{},
		model.Resource{},
	)

	// 初始化一个admin用户
	service.UserService.Add(&model.User{
		Username: "admin",
		Password: "Admin123!@#",
		RealName: "管理员",
		Status:   1,
	})

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

func main1() {
	database.Initial()
	defer database.Close()

	database.DB.Sync2(
		model.User{},
		model.Role{},
		model.UserExtend{},
		model.UserRole{},
		model.Resource{},
	)

	r, err := controller.UserController.Login([]byte("{\"username\":\"admin\",\"password\":\"Admin123!@#\"}"))
	if err != nil {
		logger.Errorln(err)
	} else {
		logger.Debugln(r)
	}
}
