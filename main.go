package main

import (
	"errors"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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

	_ = database.AutoMigrate(
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
		if err := database.DB.First(adminUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				adminUser.ID = util.SnowflakeId()
				adminUser.RealName = "管理员"
				adminUser.Status = 1
				pwd, _ := bcrypt.GenerateFromPassword([]byte("Admin123!@#"), bcrypt.DefaultCost)
				adminUser.Password = string(pwd)
				_ = database.DB.Create(adminUser)
			} else {
				logger.Errorln(err)
			}
		}
	}

	// 初始化一个超级管理员角色
	superAdminRole := &model.Role{
		RoleType: -1,
	}
	{
		if err := database.DB.First(superAdminRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				superAdminRole.ID = util.SnowflakeId()
				superAdminRole.RoleName = "超级管理员"
				_ = database.DB.Create(superAdminRole)
			} else {
				logger.Errorln(err)
			}
		}
	}

	// 关联admin和超级管理员角色
	{
		relation := &model.UserRole{UserID: adminUser.ID, RoleID: superAdminRole.ID}
		if err := database.DB.Take(relation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				relation.ID = util.SnowflakeId()
				_ = database.DB.Create(relation)
			} else {
				logger.Errorln(err)
			}
		}
	}

	return nil
}

func (UC) InjectCall(code string, headers map[string][]string, value []byte) ([]byte, error) {
	return controller.Dispatch(code, headers, value)
}

func init() {
	config.Set("moduleName", constant.ModuleName)
	config.Set("logger.level", "info")
	config.InitialLogger()
	util.InitNode(1)
}

func main() {
	defer database.Close()
	shared.ModuleServe(constant.ModuleName, &UC{})
}
