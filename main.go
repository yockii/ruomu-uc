package main

import (
	"encoding/json"
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"

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
		model.Role{},
		model.UserExtend{},
		model.UserRole{},
		model.Resource{},
	)

	// 初始化一个admin用户
	adminUser := &model.User{
		Username: "admin",
	}
	{
		if exists, err := database.DB.Get(adminUser); err != nil {
			logger.Errorln(err)
		} else if !exists {
			adminUser.Id = util.SnowflakeId()
			adminUser.RealName = "管理员"
			adminUser.Status = 1
			pwd, _ := bcrypt.GenerateFromPassword([]byte("Admin123!@#"), bcrypt.DefaultCost)
			adminUser.Password = string(pwd)
			_, _ = database.DB.Insert(adminUser)
		}
	}

	// 初始化一个超级管理员角色
	superAdminRole := &model.Role{
		RoleType: 99,
	}
	{
		if exists, err := database.DB.Get(superAdminRole); err != nil {
			logger.Errorln(err)
		} else if !exists {
			superAdminRole.Id = util.SnowflakeId()
			superAdminRole.RoleName = "超级管理员"
			_, _ = database.DB.Insert(superAdminRole)
		}
	}

	// 关联admin和超级管理员角色
	{
		relation := &model.UserRole{UserId: adminUser.Id, RoleId: superAdminRole.Id}
		if exists, err := database.DB.Get(relation); err != nil {
			logger.Errorln(err)
		} else if !exists {
			relation.Id = util.SnowflakeId()
			_, _ = database.DB.Insert(relation)
		}
	}

	return nil
}

func (UC) InjectCall(code string, headers map[string]string, value []byte) ([]byte, error) {
	return controller.Dispatch(code, headers, value)
}

func init() {
	config.Set("moduleName", constant.ModuleName)
	config.Set("logger.level", "debug")
	config.InitialLogger()
	util.InitNode(1)
}

func main2() {
	s := `{
    "id": "1600871296881659904",
    "realName": "1111",
    "status": 2
}`
	u := new(model.User)
	err := json.Unmarshal([]byte(s), u)
	fmt.Println(err, u)
}

func main() {
	defer database.Close()
	shared.ModuleServe(constant.ModuleName, &UC{})
}

func main0() {
	database.Initial()
	defer database.Close()

	database.DB.Sync2(
		model.User{},
		model.Role{},
		model.UserExtend{},
		model.UserRole{},
		model.Resource{},
	)

	//r, err := controller.UserController.Login([]byte("{\"username\":\"admin\",\"password\":\"Admin123!@#\"}"))
	r, err := controller.UserController.Update([]byte(`{
		"id": 1600714423612215296,
		"realName": "1111",
		"status": 2
	}`))
	if err != nil {
		logger.Errorln(err)
	} else {
		logger.Debugln(r)
	}
}
