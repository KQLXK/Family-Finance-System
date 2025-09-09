package model

import (
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
)

// FamilyDao 家庭数据访问对象
type FamilyDao struct{}

var (
	familyOnce sync.Once
	familyDao  *FamilyDao
)

// NewFamilyDaoInstance 返回 FamilyDao 单例实例
func NewFamilyDaoInstance() *FamilyDao {
	familyOnce.Do(func() {
		familyDao = &FamilyDao{}
	})
	return familyDao
}

// CreateFamily 创建家庭
func (FamilyDao) CreateFamily(family *Family) error {
	if err := database.DB.Create(family).Error; err != nil {
		log.Printf("创建家庭失败: %v", err)
		return err
	}
	return nil
}

// GetFamilyByID 根据ID获取家庭
func (FamilyDao) GetFamilyByID(id uint) (*Family, error) {
	var family Family
	if err := database.DB.Preload("Members").First(&family, id).Error; err != nil {
		log.Printf("获取家庭失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &family, nil
}

// GetAllFamilies 获取所有家庭
func (FamilyDao) GetAllFamilies() ([]Family, error) {
	var families []Family
	if err := database.DB.Preload("Members").Find(&families).Error; err != nil {
		log.Printf("获取所有家庭失败: %v", err)
		return nil, err
	}
	return families, nil
}

// UpdateFamily 更新家庭信息
func (FamilyDao) UpdateFamily(family *Family) error {
	if err := database.DB.Model(family).Select("name", "updated_at").Updates(family).Error; err != nil {
		log.Printf("更新家庭失败 ID=%d: %v", family.ID, err)
		return err
	}
	return nil
}

// DeleteFamily 删除家庭
func (FamilyDao) DeleteFamily(id uint) error {
	if err := database.DB.Delete(&Family{}, id).Error; err != nil {
		log.Printf("删除家庭失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}
