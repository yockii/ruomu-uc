package main

import (
	"errors"
	"github.com/yockii/ruomu-core/shared"
	moduleModel "github.com/yockii/ruomu-module/model"
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
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
		&model.User{},
		&model.Role{},
		&model.UserExtend{},
		&model.UserRole{},
		&model.Resource{},
		&model.RoleResource{},
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
	config.Set("moduleName", constant.ModuleCode)
	config.Set("logger.level", "info")
	config.InitialLogger()
	_ = util.InitNode(1)
}

func main() {
	defer database.Close()

	// 检查是否有启动参数 --mc
	args := os.Args
	runningInMicroCore := false
	for _, arg := range args {
		if arg == "--mc" {
			runningInMicroCore = true
			break
		}
	}
	if runningInMicroCore {
		shared.ModuleServe(constant.ModuleCode, &UC{})
	} else {
		registerModule()
		logger.Info("UC模块注册完成")
	}
}

func registerModule() {
	UC{}.Initial(map[string]string{})

	// 直接写表数据即可
	m := &moduleModel.Module{
		Code: constant.ModuleCode,
	}
	database.DB.Where(&moduleModel.Module{
		Code: constant.ModuleCode,
	}).Attrs(&moduleModel.Module{
		ID:     util.SnowflakeId(),
		Name:   constant.ModuleName,
		Code:   constant.ModuleCode,
		Cmd:    "./plugins/ruomu-uc --mc",
		Status: 1,
	}).FirstOrCreate(m)

	// 注入信息
	{
		mjiList := []*moduleModel.ModuleInjectInfo{
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取用户角色ID列表",
				Type:              51,
				InjectCode:        "authorizationInfoByUserId",
				AuthorizationCode: "inner",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取角色资源列表",
				Type:              51,
				InjectCode:        "authorizationInfoByRoleId",
				AuthorizationCode: "inner",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "用户登录",
				Type:              2,
				InjectCode:        constant.InjectCodeUserLogin,
				AuthorizationCode: "anon",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增用户",
				Type:              2,
				InjectCode:        constant.InjectCodeUserAdd,
				AuthorizationCode: "user:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改用户",
				Type:              3,
				InjectCode:        constant.InjectCodeUserDelete,
				AuthorizationCode: "user:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除用户",
				Type:              4,
				InjectCode:        constant.InjectCodeUserDelete,
				AuthorizationCode: "user:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取用户列表",
				Type:              1,
				InjectCode:        constant.InjectCodeUserList,
				AuthorizationCode: "user:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取用户信息",
				Type:              1,
				InjectCode:        constant.InjectCodeUserInstance,
				AuthorizationCode: "user:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改用户密码",
				Type:              3,
				InjectCode:        constant.InjectCodeUserPassword,
				AuthorizationCode: "user:password",
			},
		}
		for _, mji := range mjiList {
			temp := new(moduleModel.ModuleInjectInfo)
			database.DB.Where(&moduleModel.ModuleInjectInfo{
				ModuleID:   mji.ModuleID,
				InjectCode: mji.InjectCode,
			}).Attrs(mji).FirstOrCreate(temp)
		}
	}

	// 设置信息
	{
		for k, v := range config.GetStringMapString("database") {
			s := &moduleModel.ModuleSettings{
				ID:       util.SnowflakeId(),
				ModuleID: m.ID,
				Code:     "database." + k,
				Value:    v,
			}
			t := new(moduleModel.ModuleSettings)
			database.DB.Where(&moduleModel.ModuleSettings{
				ModuleID: m.ID,
				Code:     s.Code,
			}).Attrs(s).FirstOrCreate(t)
		}
		s := &moduleModel.ModuleSettings{
			ID:       util.SnowflakeId(),
			ModuleID: m.ID,
			Code:     "userTokenExpire",
			Value:    config.GetString("userTokenExpire"),
		}
		t := new(moduleModel.ModuleSettings)
		database.DB.Where(&moduleModel.ModuleSettings{
			ModuleID: m.ID,
			Code:     "userTokenExpire",
		}).Attrs(s).FirstOrCreate(t)
	}
}
