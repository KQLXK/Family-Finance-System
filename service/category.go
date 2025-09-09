// service/category_service.go
package service

import (
	"errors"
	"fmt"
	"github.com/KQLXK/Family-Finance-System/model"
	"strconv"
	"strings"
)

// CategoryService 分类服务接口
type CategoryService interface {
	CreateCategory(category *model.Category) error
	GetCategoryByID(id uint) (*model.Category, error)
	GetCategoriesByType(categoryType model.CategoryType) ([]model.Category, error)
	GetCategoryTreeByType(categoryType model.CategoryType) ([]model.Category, error)
	GetCategoriesByParentID(parentID uint) ([]model.Category, error)
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(id uint) error
	CategoryExists(id uint) (bool, error)
	GetFullCategoryPath(id uint) (string, error)
}

// categoryService 分类服务实现
type categoryService struct {
	categoryDao model.CategoryDao
}

// NewCategoryService 创建分类服务实例
func NewCategoryService() CategoryService {
	return &categoryService{
		categoryDao: *model.NewCategoryDaoInstance(),
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(category *model.Category) error {
	// 验证分类数据
	if err := s.validateCategory(category); err != nil {
		return err
	}

	// 检查分类名称是否已存在（同一父分类下）
	exists, err := s.categoryNameExists(category.Name, category.Type, category.ParentID, 0)
	if err != nil {
		return fmt.Errorf("检查分类名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("分类名称已存在")
	}

	// 如果有父分类，验证父分类是否存在且类型一致
	if category.ParentID != nil && *category.ParentID != 0 {
		parentCategory, err := s.categoryDao.GetCategoryByID(*category.ParentID)
		if err != nil {
			return fmt.Errorf("获取父分类失败: %v", err)
		}
		if parentCategory == nil {
			return errors.New("父分类不存在")
		}
		if parentCategory.Type != category.Type {
			return errors.New("父分类类型与子分类类型不一致")
		}
		if parentCategory.IsDeleted {
			return errors.New("父分类已被删除")
		}

		// 设置层级和路径
		category.Level = parentCategory.Level + 1
		category.Path = fmt.Sprintf("%s/%d", parentCategory.Path, parentCategory.ID)
	} else {
		// 根分类
		category.Level = 1
		category.Path = "/"
	}

	// 设置默认值
	category.IsDeleted = false

	// 创建分类
	if err := s.categoryDao.CreateCategory(category); err != nil {
		return fmt.Errorf("创建分类失败: %v", err)
	}

	// 更新路径，包含自己的ID
	category.Path = fmt.Sprintf("%s/%d", category.Path, category.ID)
	if err := s.categoryDao.UpdateCategory(category); err != nil {
		return fmt.Errorf("更新分类路径失败: %v", err)
	}

	return nil
}

// GetCategoryByID 根据ID获取分类
func (s *categoryService) GetCategoryByID(id uint) (*model.Category, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的分类ID")
	}

	// 获取分类
	category, err := s.categoryDao.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取分类失败: %v", err)
	}

	if category == nil || category.IsDeleted {
		return nil, errors.New("分类不存在或已被删除")
	}

	return category, nil
}

// GetCategoriesByType 根据类型获取分类列表
func (s *categoryService) GetCategoriesByType(categoryType model.CategoryType) ([]model.Category, error) {
	// 验证类型
	if !s.isValidCategoryType(categoryType) {
		return nil, errors.New("无效的分类类型")
	}

	// 获取分类列表
	categories, err := s.categoryDao.GetCategoriesByType(categoryType)
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 过滤已删除的分类
	var result []model.Category
	for _, category := range categories {
		if !category.IsDeleted {
			result = append(result, category)
		}
	}

	return result, nil
}

// GetCategoryTreeByType 根据类型获取分类树
func (s *categoryService) GetCategoryTreeByType(categoryType model.CategoryType) ([]model.Category, error) {
	// 验证类型
	if !s.isValidCategoryType(categoryType) {
		return nil, errors.New("无效的分类类型")
	}

	// 获取所有该类型的分类
	categories, err := s.categoryDao.GetCategoriesByType(categoryType)
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 构建分类树
	return s.buildCategoryTree(categories), nil
}

// GetCategoriesByParentID 根据父分类ID获取子分类列表
func (s *categoryService) GetCategoriesByParentID(parentID uint) ([]model.Category, error) {
	// 如果有父分类，验证父分类是否存在
	if parentID != 0 {
		parentCategory, err := s.categoryDao.GetCategoryByID(parentID)
		if err != nil {
			return nil, fmt.Errorf("获取父分类失败: %v", err)
		}
		if parentCategory == nil || parentCategory.IsDeleted {
			return nil, errors.New("父分类不存在或已被删除")
		}
	}

	// 获取子分类列表
	categories, err := s.categoryDao.GetCategoriesByParentID(parentID)
	if err != nil {
		return nil, fmt.Errorf("获取子分类列表失败: %v", err)
	}

	// 过滤已删除的分类
	var result []model.Category
	for _, category := range categories {
		if !category.IsDeleted {
			result = append(result, category)
		}
	}

	return result, nil
}

// GetAllCategories 获取所有分类
func (s *categoryService) GetAllCategories() ([]model.Category, error) {
	// 获取所有分类
	categories, err := s.categoryDao.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("获取所有分类失败: %v", err)
	}

	// 过滤已删除的分类
	var result []model.Category
	for _, category := range categories {
		if !category.IsDeleted {
			result = append(result, category)
		}
	}

	return result, nil
}

// UpdateCategory 更新分类信息
func (s *categoryService) UpdateCategory(category *model.Category) error {
	// 验证分类ID
	if category.ID == 0 {
		return errors.New("无效的分类ID")
	}

	// 验证分类数据
	if err := s.validateCategory(category); err != nil {
		return err
	}

	// 检查分类是否存在
	exists, err := s.CategoryExists(category.ID)
	if err != nil {
		return fmt.Errorf("检查分类是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("分类不存在")
	}

	// 检查分类名称是否已被其他分类使用（同一父分类下）
	exists, err = s.categoryNameExists(category.Name, category.Type, category.ParentID, category.ID)
	if err != nil {
		return fmt.Errorf("检查分类名称是否已存在时出错: %v", err)
	}
	if exists {
		return errors.New("分类名称已存在")
	}

	// 如果有父分类，验证父分类是否存在且类型一致
	if category.ParentID != nil && *category.ParentID != 0 {
		parentCategory, err := s.categoryDao.GetCategoryByID(*category.ParentID)
		if err != nil {
			return fmt.Errorf("获取父分类失败: %v", err)
		}
		if parentCategory == nil || parentCategory.IsDeleted {
			return errors.New("父分类不存在或已被删除")
		}
		if parentCategory.Type != category.Type {
			return errors.New("父分类类型与子分类类型不一致")
		}

		// 检查是否形成了循环依赖
		if err := s.checkCircularDependency(category.ID, *category.ParentID); err != nil {
			return err
		}

		// 更新层级和路径
		category.Level = parentCategory.Level + 1
		category.Path = fmt.Sprintf("%s/%d", parentCategory.Path, category.ID)
	} else {
		// 根分类
		category.Level = 1
		category.Path = fmt.Sprintf("/%d", category.ID)
	}

	// 更新分类信息
	if err := s.categoryDao.UpdateCategory(category); err != nil {
		return fmt.Errorf("更新分类信息失败: %v", err)
	}

	return nil
}

// DeleteCategory 删除分类（软删除）
func (s *categoryService) DeleteCategory(id uint) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的分类ID")
	}

	// 检查分类是否存在
	category, err := s.categoryDao.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("检查分类是否存在时出错: %v", err)
	}
	if category == nil || category.IsDeleted {
		return errors.New("分类不存在或已被删除")
	}

	// 检查是否有子分类
	children, err := s.categoryDao.GetCategoriesByParentID(id)
	if err != nil {
		return fmt.Errorf("检查子分类时出错: %v", err)
	}

	// 过滤未删除的子分类
	var activeChildren []model.Category
	for _, child := range children {
		if !child.IsDeleted {
			activeChildren = append(activeChildren, child)
		}
	}

	if len(activeChildren) > 0 {
		return errors.New("无法删除包含子分类的分类，请先删除所有子分类")
	}

	// 软删除分类（设置IsDeleted为true）
	if err := s.categoryDao.DeleteCategory(id); err != nil {
		return fmt.Errorf("删除分类失败: %v", err)
	}

	return nil
}

// CategoryExists 检查分类是否存在
func (s *categoryService) CategoryExists(id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}

	category, err := s.categoryDao.GetCategoryByID(id)
	if err != nil {
		return false, err
	}

	return category != nil && !category.IsDeleted, nil
}

// GetFullCategoryPath 获取完整分类路径
func (s *categoryService) GetFullCategoryPath(id uint) (string, error) {
	// 验证ID
	if id == 0 {
		return "", errors.New("无效的分类ID")
	}

	// 获取分类
	category, err := s.categoryDao.GetCategoryByID(id)
	if err != nil {
		return "", fmt.Errorf("获取分类失败: %v", err)
	}

	if category == nil || category.IsDeleted {
		return "", errors.New("分类不存在或已被删除")
	}

	// 解析路径中的分类ID
	pathIDs := strings.Split(strings.Trim(category.Path, "/"), "/")
	var pathNames []string

	// 获取每个分类的名称
	for _, pathIDStr := range pathIDs {
		if pathIDStr == "" {
			continue
		}

		pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
		if err != nil {
			return "", fmt.Errorf("解析路径ID失败: %v", err)
		}

		pathCategory, err := s.categoryDao.GetCategoryByID(uint(pathID))
		if err != nil {
			return "", fmt.Errorf("获取路径分类失败: %v", err)
		}

		if pathCategory != nil && !pathCategory.IsDeleted {
			pathNames = append(pathNames, pathCategory.Name)
		}
	}

	return strings.Join(pathNames, " > "), nil
}

// validateCategory 验证分类数据
func (s *categoryService) validateCategory(category *model.Category) error {
	// 验证名称
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("分类名称不能为空")
	}

	if len(category.Name) > 100 {
		return errors.New("分类名称长度不能超过100个字符")
	}

	// 验证类型
	if !s.isValidCategoryType(category.Type) {
		return errors.New("无效的分类类型")
	}

	return nil
}

// categoryNameExists 检查分类名称是否已存在
func (s *categoryService) categoryNameExists(name string, categoryType model.CategoryType, parentID *uint, excludeID uint) (bool, error) {
	categories, err := s.categoryDao.GetAllCategories()
	if err != nil {
		return false, err
	}

	for _, category := range categories {
		if category.ID != excludeID &&
			!category.IsDeleted &&
			category.Type == categoryType &&
			strings.EqualFold(category.Name, name) {

			// 检查父分类是否相同
			if (parentID == nil && category.ParentID == nil) ||
				(parentID != nil && category.ParentID != nil && *parentID == *category.ParentID) {
				return true, nil
			}
		}
	}

	return false, nil
}

// isValidCategoryType 验证分类类型是否有效
func (s *categoryService) isValidCategoryType(categoryType model.CategoryType) bool {
	switch categoryType {
	case model.CategoryIncome, model.CategoryExpense:
		return true
	default:
		return false
	}
}

// buildCategoryTree 构建分类树
func (s *categoryService) buildCategoryTree(categories []model.Category) []model.Category {
	// 创建ID到分类的映射
	categoryMap := make(map[uint]*model.Category)
	for i := range categories {
		if !categories[i].IsDeleted {
			categoryMap[categories[i].ID] = &categories[i]
		}
	}

	// 构建树结构
	var roots []model.Category
	for i := range categories {
		if categories[i].IsDeleted {
			continue
		}

		if categories[i].ParentID == nil || *categories[i].ParentID == 0 {
			// 根分类
			roots = append(roots, categories[i])
		} else if parent, exists := categoryMap[*categories[i].ParentID]; exists {
			// 子分类，添加到父分类的Children中
			parent.Children = append(parent.Children, categories[i])
		}
	}

	return roots
}

// checkCircularDependency 检查循环依赖
func (s *categoryService) checkCircularDependency(categoryID uint, parentID uint) error {
	// 检查是否试图将分类设置为自己的父分类
	if categoryID == parentID {
		return errors.New("不能将分类设置为自己的父分类")
	}

	// 检查是否形成了循环依赖
	currentParentID := parentID
	for currentParentID != 0 {
		parentCategory, err := s.categoryDao.GetCategoryByID(currentParentID)
		if err != nil {
			return fmt.Errorf("检查循环依赖时出错: %v", err)
		}

		if parentCategory == nil || parentCategory.IsDeleted {
			break
		}

		// 如果找到了循环依赖
		if parentCategory.ID == categoryID {
			return errors.New("检测到循环依赖，不能将子分类设置为父分类")
		}

		// 继续向上检查
		if parentCategory.ParentID == nil {
			break
		}
		currentParentID = *parentCategory.ParentID
	}

	return nil
}
