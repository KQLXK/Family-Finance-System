package model

import (
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
)

// TagDao 标签数据访问对象
type TagDao struct{}

var (
	tagOnce sync.Once
	tagDao  *TagDao
)

// NewTagDaoInstance 返回 TagDao 单例实例
func NewTagDaoInstance() *TagDao {
	tagOnce.Do(func() {
		tagDao = &TagDao{}
	})
	return tagDao
}

// CreateTag 创建标签
func (TagDao) CreateTag(tag *Tag) error {
	if err := database.DB.Create(tag).Error; err != nil {
		log.Printf("创建标签失败: %v", err)
		return err
	}
	return nil
}

// GetTagByID 根据ID获取标签
func (TagDao) GetTagByID(id uint) (*Tag, error) {
	var tag Tag
	if err := database.DB.Preload("Family").First(&tag, id).Error; err != nil {
		log.Printf("获取标签失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &tag, nil
}

// GetTagsByFamilyID 根据家庭ID获取标签列表
func (TagDao) GetTagsByFamilyID(familyID uint) ([]Tag, error) {
	var tags []Tag
	if err := database.DB.Where("family_id = ? AND is_active = true", familyID).Find(&tags).Error; err != nil {
		log.Printf("获取家庭标签失败 FamilyID=%d: %v", familyID, err)
		return nil, err
	}
	return tags, nil
}

// GetTagsByType 根据类型获取标签列表
func (TagDao) GetTagsByType(familyID uint, tagType string) ([]Tag, error) {
	var tags []Tag
	if err := database.DB.Where("family_id = ? AND type = ? AND is_active = true", familyID, tagType).Find(&tags).Error; err != nil {
		log.Printf("获取类型标签失败 FamilyID=%d, Type=%s: %v", familyID, tagType, err)
		return nil, err
	}
	return tags, nil
}

// UpdateTag 更新标签信息
func (TagDao) UpdateTag(tag *Tag) error {
	if err := database.DB.Save(tag).Error; err != nil {
		log.Printf("更新标签失败 ID=%d: %v", tag.ID, err)
		return err
	}
	return nil
}

// DeleteTag 软删除标签
func (TagDao) DeleteTag(id uint) error {
	if err := database.DB.Model(&Tag{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		log.Printf("删除标签失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}

func (TagDao) GetAllTags() ([]Tag, error) {
	var tags []Tag
	if err := database.DB.Preload("Family").Find(&tags).Error; err != nil {
		log.Printf("获取所有标签失败: %v", err)
		return nil, err
	}
	return tags, nil
}

// IsTagUsedInTransactions 检查标签是否被交易使用
func (TagDao) IsTagUsedInTransactions(id uint) (bool, error) {
	var count int64
	if err := database.DB.Model(&TransactionTag{}).
		Where("tag_id = ?", id).
		Count(&count).Error; err != nil {
		log.Printf("检查标签是否被使用失败 TagID=%d: %v", id, err)
		return false, err
	}
	return count > 0, nil
}
