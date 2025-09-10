// service/tag_service.go
package service

import (
	"errors"
	"fmt"
	"github.com/KQLXK/Family-Finance-System/model"
	"strings"
)

// TagService 标签服务接口
type TagService interface {
	CreateTag(tag *model.Tag) error
	GetTagByID(id uint) (*model.Tag, error)
	GetTagsByFamilyID(familyID uint) ([]model.Tag, error)
	GetTagsByType(familyID uint, tagType string) ([]model.Tag, error)
	GetAllTags() ([]model.Tag, error)
	UpdateTag(tag *model.Tag) error
	DeleteTag(id uint) error
	TagExists(id uint) (bool, error)
}

// tagService 标签服务实现
type tagService struct {
	tagDao    model.TagDao
	familyDao model.FamilyDao
}

// NewTagService 创建标签服务实例
func NewTagService() TagService {
	return &tagService{
		tagDao:    *model.NewTagDaoInstance(),
		familyDao: *model.NewFamilyDaoInstance(),
	}
}

// CreateTag 创建标签
func (s *tagService) CreateTag(tag *model.Tag) error {
	// 验证标签数据
	if err := s.validateTag(tag); err != nil {
		return err
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(tag.FamilyID)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return errors.New("关联的家庭不存在")
	}

	// 检查标签名称是否已存在（同一家庭下）
	exists, err := s.tagNameExists(tag.Name, tag.FamilyID, 0)
	if err != nil {
		return fmt.Errorf("检查标签名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("标签名称已存在")
	}

	// 设置默认值
	tag.IsActive = true

	// 创建标签
	if err := s.tagDao.CreateTag(tag); err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}

	return nil
}

// GetTagByID 根据ID获取标签
func (s *tagService) GetTagByID(id uint) (*model.Tag, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的标签ID")
	}

	// 获取标签
	tag, err := s.tagDao.GetTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}

	if tag == nil || !tag.IsActive {
		return nil, errors.New("标签不存在或已被禁用")
	}

	return tag, nil
}

// GetTagsByFamilyID 根据家庭ID获取标签列表
func (s *tagService) GetTagsByFamilyID(familyID uint) ([]model.Tag, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 获取标签列表
	tags, err := s.tagDao.GetTagsByFamilyID(familyID)
	if err != nil {
		return nil, fmt.Errorf("获取家庭标签失败: %v", err)
	}

	// 过滤禁用的标签
	var result []model.Tag
	for _, tag := range tags {
		if tag.IsActive {
			result = append(result, tag)
		}
	}

	return result, nil
}

// GetTagsByType 根据类型获取标签列表
func (s *tagService) GetTagsByType(familyID uint, tagType string) ([]model.Tag, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 获取标签列表
	tags, err := s.tagDao.GetTagsByType(familyID, tagType)
	if err != nil {
		return nil, fmt.Errorf("获取类型标签失败: %v", err)
	}

	// 过滤禁用的标签
	var result []model.Tag
	for _, tag := range tags {
		if tag.IsActive {
			result = append(result, tag)
		}
	}

	return result, nil
}

// GetAllTags 获取所有标签
func (s *tagService) GetAllTags() ([]model.Tag, error) {
	// 获取所有标签
	tags, err := s.tagDao.GetAllTags()
	if err != nil {
		return nil, fmt.Errorf("获取所有标签失败: %v", err)
	}

	// 过滤禁用的标签
	var result []model.Tag
	for _, tag := range tags {
		if tag.IsActive {
			result = append(result, tag)
		}
	}

	return result, nil
}

// UpdateTag 更新标签信息
func (s *tagService) UpdateTag(tag *model.Tag) error {
	// 验证标签ID
	if tag.ID == 0 {
		return errors.New("无效的标签ID")
	}

	// 验证标签数据
	if err := s.validateTag(tag); err != nil {
		return err
	}

	// 检查标签是否存在
	exists, err := s.TagExists(tag.ID)
	if err != nil {
		return fmt.Errorf("检查标签是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("标签不存在")
	}

	// 检查标签名称是否已被其他标签使用（同一家庭下）
	exists, err = s.tagNameExists(tag.Name, tag.FamilyID, tag.ID)
	if err != nil {
		return fmt.Errorf("检查标签名称是否已存在时出错: %v", err)
	}
	if exists {
		return errors.New("标签名称已存在")
	}

	// 更新标签信息
	if err := s.tagDao.UpdateTag(tag); err != nil {
		return fmt.Errorf("更新标签信息失败: %v", err)
	}

	return nil
}

// DeleteTag 删除标签（软删除，设置IsActive为false）
func (s *tagService) DeleteTag(id uint) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的标签ID")
	}

	// 检查标签是否存在
	exists, err := s.TagExists(id)
	if err != nil {
		return fmt.Errorf("检查标签是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("标签不存在")
	}

	// 检查标签是否被交易使用
	isUsed, err := s.tagDao.IsTagUsedInTransactions(id)
	if err != nil {
		return fmt.Errorf("检查标签是否被使用时出错: %v", err)
	}
	if isUsed {
		return errors.New("标签已被交易使用，无法删除")
	}

	// 软删除标签（设置IsActive为false）
	if err := s.tagDao.DeleteTag(id); err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}

	return nil
}

// TagExists 检查标签是否存在
func (s *tagService) TagExists(id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}

	tag, err := s.tagDao.GetTagByID(id)
	if err != nil {
		return false, err
	}

	return tag != nil && tag.IsActive, nil
}

// validateTag 验证标签数据
func (s *tagService) validateTag(tag *model.Tag) error {
	// 验证名称
	if strings.TrimSpace(tag.Name) == "" {
		return errors.New("标签名称不能为空")
	}

	if len(tag.Name) > 100 {
		return errors.New("标签名称长度不能超过100个字符")
	}

	// 验证类型
	if strings.TrimSpace(tag.Type) == "" {
		return errors.New("标签类型不能为空")
	}

	if len(tag.Type) > 50 {
		return errors.New("标签类型长度不能超过50个字符")
	}

	// 验证颜色格式（如果提供了颜色）
	if tag.Color != "" {
		if !strings.HasPrefix(tag.Color, "#") || len(tag.Color) != 7 {
			return errors.New("颜色格式不正确，应为#开头加6位十六进制数")
		}
	}

	return nil
}

// tagNameExists 检查标签名称是否已存在
func (s *tagService) tagNameExists(name string, familyID uint, excludeID uint) (bool, error) {
	tags, err := s.tagDao.GetAllTags()
	if err != nil {
		return false, err
	}

	for _, tag := range tags {
		if tag.ID != excludeID &&
			tag.IsActive &&
			tag.FamilyID == familyID &&
			strings.EqualFold(tag.Name, name) {
			return true, nil
		}
	}

	return false, nil
}

// familyExists 检查家庭是否存在
func (s *tagService) familyExists(familyID uint) (bool, error) {
	if familyID == 0 {
		return false, nil
	}

	family, err := s.familyDao.GetFamilyByID(familyID)
	if err != nil {
		return false, err
	}

	return family != nil, nil
}
