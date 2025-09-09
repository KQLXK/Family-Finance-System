package model

import (
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
)

// CategoryDao 分类数据访问对象
type CategoryDao struct{}

var (
	categoryOnce sync.Once
	categoryDao  *CategoryDao
)

// NewCategoryDaoInstance 返回 CategoryDao 单例实例
func NewCategoryDaoInstance() *CategoryDao {
	categoryOnce.Do(func() {
		categoryDao = &CategoryDao{}
	})
	return categoryDao
}

// CreateCategory 创建分类
func (CategoryDao) CreateCategory(category *Category) error {
	if err := database.DB.Create(category).Error; err != nil {
		log.Printf("创建分类失败: %v", err)
		return err
	}
	return nil
}

// GetCategoryByID 根据ID获取分类
func (CategoryDao) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	if err := database.DB.Preload("Parent").Preload("Children").First(&category, id).Error; err != nil {
		log.Printf("获取分类失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &category, nil
}

// GetCategoriesByType 根据类型获取分类列表
func (CategoryDao) GetCategoriesByType(categoryType CategoryType) ([]Category, error) {
	var categories []Category
	if err := database.DB.Where("type = ? AND is_deleted = false", categoryType).
		Preload("Children", "is_deleted = false").
		Find(&categories).Error; err != nil {
		log.Printf("获取分类列表失败 Type=%s: %v", categoryType, err)
		return nil, err
	}
	return categories, nil
}

// GetCategoriesByParentID 根据父分类ID获取子分类列表
func (CategoryDao) GetCategoriesByParentID(parentID uint) ([]Category, error) {
	var categories []Category
	if err := database.DB.Where("parent_id = ? AND is_deleted = false", parentID).Find(&categories).Error; err != nil {
		log.Printf("获取子分类失败 ParentID=%d: %v", parentID, err)
		return nil, err
	}
	return categories, nil
}

// UpdateCategory 更新分类信息
func (CategoryDao) UpdateCategory(category *Category) error {
	if err := database.DB.Save(category).Error; err != nil {
		log.Printf("更新分类失败 ID=%d: %v", category.ID, err)
		return err
	}
	return nil
}

// DeleteCategory 软删除分类
func (CategoryDao) DeleteCategory(id uint) error {
	if err := database.DB.Model(&Category{}).Where("id = ?", id).Update("is_deleted", true).Error; err != nil {
		log.Printf("删除分类失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}

func (CategoryDao) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := database.DB.Preload("Parent").Preload("Children").Find(&categories).Error; err != nil {
		log.Printf("获取所有分类失败: %v", err)
		return nil, err
	}
	return categories, nil
}
