package service

import (
	"errors"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"

	"github.com/yockii/ruomu-uc/model"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) Add(instance *model.User) (duplicated bool, success bool, err error) {
	if instance.Username == "" {
		err = errors.New("username is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.User{}).Where(&model.User{Username: instance.Username, ExternalType: instance.ExternalType}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()
	if instance.Password != "" {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(instance.Password), bcrypt.DefaultCost)
		instance.Password = string(pwd)
	}
	instance.Status = 1
	err = database.DB.Create(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 完成后密码置空
	instance.Password = ""
	success = true
	return
}

func (s *userService) PaginateBetweenTimes(condition *model.User, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.User, err error) {
	// 处理不允许查询的字段
	if condition.Password != "" {
		condition.Password = ""
	}
	tx := database.DB.Limit(100)
	if limit > -1 {
		tx.Limit(limit)
	}
	if offset > -1 {
		tx.Offset(offset)
	}

	if orderBy != "" {
		tx.Order(orderBy)
	}

	// 处理时间字段，在某段时间之间
	for tc, tr := range tcList {
		if tc != "" {
			if !tr.Start.IsZero() && !tr.End.IsZero() {
				tx.Where(tc+" between ? and ?", time.Time(tr.Start), time.Time(tr.End))
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx.Where(tc+" <= ?", time.Time(tr.End))
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx.Where(tc+" > ?", time.Time(tr.Start))
			}
		}
	}

	// 模糊查找
	if condition.Username != "" {
		tx.Where("username like ?", condition.Username+"%")
		condition.Username = ""
	}
	err = tx.Omit("password").Find(&list, condition).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	return total, list, nil
}
