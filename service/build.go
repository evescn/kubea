package service

import (
	"errors"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Build build

type build struct{}

// List 列表
func (*build) List(repoName string, page, limit int) (*dao.Builds, error) {
	return dao.Build.List(repoName, page, limit)
}

// Add 新增
func (*build) Add(b *model.Build) error {
	_, has, err := dao.Build.Has(b.RepoName, b.Tag)
	if err != nil {
		return err
	}

	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	return nil
}

// Update 更新
func (*build) Update(b *model.Build) error {
	return dao.Build.Update(b)
}
