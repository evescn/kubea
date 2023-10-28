package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea/db"
	"kubea/model"
)

var Service service

type service struct{}

type ServiceRes struct {
	model.Service
	PName    string `json:"p_name"`
	Password string `json:"p_password"`
}

type Services struct {
	Items []*ServiceRes `json:"items"`
	Total int           `json:"total"`
}

// List 列表
// serviceName用于模糊查询，过滤
// page，limit用于分页
func (*service) List(serviceName string, eid uint, page, limit int) (*Services, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		serviceList = make([]*ServiceRes, 0)
		total       = 0
	)

	//数据库查询，先查total
	query := db.GORM.Model(model.Service{}).
		Select("service.*, password.p_name, password.password").
		Joins("left join password on service.pid = password.id").
		Joins("left join env on service.eid = env.id").
		Where("eid = ? ", eid)
	if serviceName != "" {
		query = query.Where("name like ?", "%"+serviceName+"%")
	}
	tx := query.Count(&total)

	if tx.Error != nil {
		logger.Error("获取Service列表失败," + tx.Error.Error())
		return nil, errors.New("获取Service列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = query.Limit(limit).
		Offset(startSet).
		Order("service.id").
		Find(&serviceList)
	if tx.Error != nil {
		logger.Error("获取Service列表失败," + tx.Error.Error())
		return nil, errors.New("获取Service列表失败," + tx.Error.Error())
	}

	return &Services{
		Items: serviceList,
		Total: total,
	}, nil
}

func (*service) ListNoUserInfo(serviceName string, eid uint, page, limit int) (*Services, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		serviceList = make([]*ServiceRes, 0)
		total       = 0
	)

	//数据库查询，先查total
	query := db.GORM.Model(model.Service{}).
		Select("service.*").
		Joins("left join env on service.eid = env.id").
		Where("eid = ? ", eid)

	if serviceName != "" {
		query = query.Where("name like ?", "%"+serviceName+"%")
	}
	tx := query.Count(&total)

	if tx.Error != nil {
		logger.Error("获取Service列表失败," + tx.Error.Error())
		return nil, errors.New("获取Service列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = query.Limit(limit).
		Offset(startSet).
		Order("service.id").
		Find(&serviceList)
	if tx.Error != nil {
		logger.Error("获取Service列表失败," + tx.Error.Error())
		return nil, errors.New("获取Service列表失败," + tx.Error.Error())
	}

	return &Services{
		Items: serviceList,
		Total: total,
	}, nil
}

// Get 根据环境名查询服务信息，用于判断环境是否还关联 URL 信息，删除环境
func (*service) Get(eid uint) (*model.Service, bool, error) {
	data := new(model.Service)
	tx := db.GORM.Where("eid = ? ", eid).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据名称查询Service失败," + tx.Error.Error())
		return nil, false, errors.New("根据名称查询Service失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据名称查询，用于代码层去重，查询账号信息
func (*service) Has(serviceName string, eid uint) (*model.Service, bool, error) {
	data := new(model.Service)
	tx := db.GORM.Where("name = ? and eid = ? ", serviceName, eid).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据名称查询Service失败," + tx.Error.Error())
		return nil, false, errors.New("根据名称查询Service失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*service) Add(e *model.Service) error {
	tx := db.GORM.Create(&e)
	if tx.Error != nil {
		logger.Error("新增Service信息失败," + tx.Error.Error())
		return errors.New("新增Service信息失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*service) Update(e *model.Service) error {
	tx := db.GORM.Model(&model.Service{}).Where("id = ?", e.ID).Updates(&e)
	if tx.Error != nil {
		logger.Error("更新Service信息失败," + tx.Error.Error())
		return errors.New("更新Service信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*service) Delete(id uint) error {
	data := new(model.Service)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Service信息失败," + tx.Error.Error())
		return errors.New("删除Service信息失败," + tx.Error.Error())
	}

	return nil
}
