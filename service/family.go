// service/family_service.go
package service

import (
	"errors"
	"fmt"
	"github.com/KQLXK/Family-Finance-System/model"
	"strings"
)

// FamilyService 家庭服务接口
type FamilyService interface {
	CreateFamily(family *model.Family) error
	GetFamilyByID(id uint) (*model.Family, error)
	GetAllFamilies() ([]model.Family, error)
	UpdateFamily(family *model.Family) error
	DeleteFamily(id uint) error
	GetFamilyWithMembers(id uint) (*model.Family, error)
	FamilyExists(id uint) (bool, error)
}

// familyService 家庭服务实现
type familyService struct {
	familyDao *model.FamilyDao
}

// NewFamilyService 创建家庭服务实例
func NewFamilyService() FamilyService {
	return &familyService{
		familyDao: model.NewFamilyDaoInstance(),
	}
}

// CreateFamily 创建家庭
func (s *familyService) CreateFamily(family *model.Family) error {
	// 验证家庭名称
	if err := s.validateFamilyName(family.Name); err != nil {
		return err
	}

	// 检查家庭名称是否已存在
	exists, err := s.familyNameExists(family.Name)
	if err != nil {
		return fmt.Errorf("检查家庭名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("家庭名称已存在")
	}

	// 创建家庭
	if err := s.familyDao.CreateFamily(family); err != nil {
		return fmt.Errorf("创建家庭失败: %v", err)
	}

	return nil
}

// GetFamilyByID 根据ID获取家庭
func (s *familyService) GetFamilyByID(id uint) (*model.Family, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 获取家庭
	family, err := s.familyDao.GetFamilyByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取家庭失败: %v", err)
	}

	if family == nil {
		return nil, errors.New("家庭不存在")
	}

	return family, nil
}

// GetAllFamilies 获取所有家庭
func (s *familyService) GetAllFamilies() ([]model.Family, error) {
	// 获取所有家庭
	families, err := s.familyDao.GetAllFamilies()
	if err != nil {
		return nil, fmt.Errorf("获取家庭列表失败: %v", err)
	}

	return families, nil
}

// UpdateFamily 更新家庭信息
func (s *familyService) UpdateFamily(family *model.Family) error {
	// 验证家庭ID
	if family.ID == 0 {
		return errors.New("无效的家庭ID")
	}

	// 验证家庭名称
	if err := s.validateFamilyName(family.Name); err != nil {
		return err
	}

	// 检查家庭是否存在
	exists, err := s.FamilyExists(family.ID)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("家庭不存在")
	}

	// 检查家庭名称是否已被其他家庭使用
	exists, err = s.otherFamilyHasSameName(family.ID, family.Name)
	if err != nil {
		return fmt.Errorf("检查家庭名称是否重复时出错: %v", err)
	}
	if exists {
		return errors.New("家庭名称已被其他家庭使用")
	}

	// 更新家庭信息
	if err := s.familyDao.UpdateFamily(family); err != nil {
		return fmt.Errorf("更新家庭信息失败: %v", err)
	}

	return nil
}

// DeleteFamily 删除家庭
func (s *familyService) DeleteFamily(id uint) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	exists, err := s.FamilyExists(id)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("家庭不存在")
	}

	// 获取家庭详情（包含成员）
	family, err := s.GetFamilyWithMembers(id)
	if err != nil {
		return fmt.Errorf("获取家庭详情失败: %v", err)
	}

	// 检查家庭是否有成员
	if len(family.Members) > 0 {
		return errors.New("无法删除包含成员的家庭，请先移除所有成员")
	}

	// 删除家庭
	if err := s.familyDao.DeleteFamily(id); err != nil {
		return fmt.Errorf("删除家庭失败: %v", err)
	}

	return nil
}

// GetFamilyWithMembers 获取家庭及其成员信息
func (s *familyService) GetFamilyWithMembers(id uint) (*model.Family, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 获取家庭及其成员
	family, err := s.familyDao.GetFamilyByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取家庭及其成员失败: %v", err)
	}

	if family == nil {
		return nil, errors.New("家庭不存在")
	}

	return family, nil
}

// FamilyExists 检查家庭是否存在
func (s *familyService) FamilyExists(id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}

	family, err := s.familyDao.GetFamilyByID(id)
	if err != nil {
		return false, err
	}

	return family != nil, nil
}

// validateFamilyName 验证家庭名称
func (s *familyService) validateFamilyName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("家庭名称不能为空")
	}

	if len(name) > 100 {
		return errors.New("家庭名称长度不能超过100个字符")
	}

	// 可以添加更多的验证规则，如特殊字符检查等

	return nil
}

// familyNameExists 检查家庭名称是否已存在
func (s *familyService) familyNameExists(name string) (bool, error) {
	families, err := s.familyDao.GetAllFamilies()
	if err != nil {
		return false, err
	}

	for _, family := range families {
		if strings.EqualFold(family.Name, name) {
			return true, nil
		}
	}

	return false, nil
}

// otherFamilyHasSameName 检查其他家庭是否使用相同的名称
func (s *familyService) otherFamilyHasSameName(id uint, name string) (bool, error) {
	families, err := s.familyDao.GetAllFamilies()
	if err != nil {
		return false, err
	}

	for _, family := range families {
		if family.ID != id && strings.EqualFold(family.Name, name) {
			return true, nil
		}
	}

	return false, nil
}
