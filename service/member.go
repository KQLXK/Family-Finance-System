// service/member_service.go
package service

import (
	"errors"
	"fmt"
	"github.com/KQLXK/Family-Finance-System/model"
	"strings"
)

// MemberService 成员服务接口
type MemberService interface {
	CreateMember(member *model.Member) error
	GetMemberByID(id uint) (*model.Member, error)
	GetMembersByFamilyID(familyID uint) ([]model.Member, error)
	GetAllMembers() ([]model.Member, error)
	UpdateMember(member *model.Member) error
	DeleteMember(id uint) error
	ChangeMemberRole(id uint, role model.MemberRole) error
	GetActiveMembersByFamilyID(familyID uint) ([]model.Member, error)
	MemberExists(id uint) (bool, error)
}

// memberService 成员服务实现
type memberService struct {
	memberDao     model.MemberDao
	familyDao     model.FamilyDao
	familyService FamilyService
}

// NewMemberService 创建成员服务实例
func NewMemberService() MemberService {
	return &memberService{
		memberDao:     *model.NewMemberDaoInstance(),
		familyDao:     *model.NewFamilyDaoInstance(),
		familyService: NewFamilyService(),
	}
}

// CreateMember 创建成员
func (s *memberService) CreateMember(member *model.Member) error {
	// 验证成员数据
	if err := s.validateMember(member); err != nil {
		return err
	}

	// 检查家庭是否存在
	familyExists, err := s.familyService.FamilyExists(member.FamilyID)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return errors.New("关联的家庭不存在")
	}

	// 检查邮箱是否已被使用（如果提供了邮箱）
	if member.Email != "" {
		emailExists, err := s.emailExists(member.Email, 0)
		if err != nil {
			return fmt.Errorf("检查邮箱是否已存在时出错: %v", err)
		}
		if emailExists {
			return errors.New("邮箱已被其他成员使用")
		}
	}

	// 检查手机号是否已被使用（如果提供了手机号）
	if member.Phone != "" {
		phoneExists, err := s.phoneExists(member.Phone, 0)
		if err != nil {
			return fmt.Errorf("检查手机号是否已存在时出错: %v", err)
		}
		if phoneExists {
			return errors.New("手机号已被其他成员使用")
		}
	}

	// 设置默认状态
	member.Status = 1

	// 创建成员
	if err := s.memberDao.CreateMember(member); err != nil {
		return fmt.Errorf("创建成员失败: %v", err)
	}

	return nil
}

// GetMemberByID 根据ID获取成员
func (s *memberService) GetMemberByID(id uint) (*model.Member, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的成员ID")
	}

	// 获取成员
	member, err := s.memberDao.GetMemberByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取成员失败: %v", err)
	}

	if member == nil {
		return nil, errors.New("成员不存在")
	}

	return member, nil
}

// GetMembersByFamilyID 根据家庭ID获取成员列表
func (s *memberService) GetMembersByFamilyID(familyID uint) ([]model.Member, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyService.FamilyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 获取成员列表
	members, err := s.memberDao.GetMembersByFamilyID(familyID)
	if err != nil {
		return nil, fmt.Errorf("获取家庭成员失败: %v", err)
	}

	return members, nil
}

// GetAllMembers 获取所有成员
func (s *memberService) GetAllMembers() ([]model.Member, error) {
	// 获取所有成员
	members, err := s.memberDao.GetAllMembers()
	if err != nil {
		return nil, fmt.Errorf("获取所有成员失败: %v", err)
	}

	return members, nil
}

// UpdateMember 更新成员信息
func (s *memberService) UpdateMember(member *model.Member) error {
	// 验证成员ID
	if member.ID == 0 {
		return errors.New("无效的成员ID")
	}

	// 验证成员数据
	if err := s.validateMember(member); err != nil {
		return err
	}

	// 检查成员是否存在
	exists, err := s.MemberExists(member.ID)
	if err != nil {
		return fmt.Errorf("检查成员是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("成员不存在")
	}

	// 检查邮箱是否已被其他成员使用
	if member.Email != "" {
		emailExists, err := s.emailExists(member.Email, member.ID)
		if err != nil {
			return fmt.Errorf("检查邮箱是否已存在时出错: %v", err)
		}
		if emailExists {
			return errors.New("邮箱已被其他成员使用")
		}
	}

	// 检查手机号是否已被其他成员使用
	if member.Phone != "" {
		phoneExists, err := s.phoneExists(member.Phone, member.ID)
		if err != nil {
			return fmt.Errorf("检查手机号是否已存在时出错: %v", err)
		}
		if phoneExists {
			return errors.New("手机号已被其他成员使用")
		}
	}

	// 更新成员信息
	if err := s.memberDao.UpdateMember(member); err != nil {
		return fmt.Errorf("更新成员信息失败: %v", err)
	}

	return nil
}

// DeleteMember 删除成员（软删除）
func (s *memberService) DeleteMember(id uint) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的成员ID")
	}

	// 检查成员是否存在
	exists, err := s.MemberExists(id)
	if err != nil {
		return fmt.Errorf("检查成员是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("成员不存在")
	}

	// 软删除成员（设置状态为0）
	if err := s.memberDao.DeleteMember(id); err != nil {
		return fmt.Errorf("删除成员失败: %v", err)
	}

	return nil
}

// ChangeMemberRole 更改成员角色
func (s *memberService) ChangeMemberRole(id uint, role model.MemberRole) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的成员ID")
	}

	// 验证角色
	if !s.isValidRole(role) {
		return errors.New("无效的成员角色")
	}

	// 检查成员是否存在
	exists, err := s.MemberExists(id)
	if err != nil {
		return fmt.Errorf("检查成员是否存在时出错: %v", err)
	}
	if !exists {
		return errors.New("成员不存在")
	}

	// 更新成员角色
	if err := s.memberDao.UpdateMemberRole(id, role); err != nil {
		return fmt.Errorf("更改成员角色失败: %v", err)
	}

	return nil
}

// GetActiveMembersByFamilyID 根据家庭ID获取活跃成员列表
func (s *memberService) GetActiveMembersByFamilyID(familyID uint) ([]model.Member, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyService.FamilyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 获取活跃成员列表
	members, err := s.memberDao.GetActiveMembersByFamilyID(familyID)
	if err != nil {
		return nil, fmt.Errorf("获取家庭活跃成员失败: %v", err)
	}

	return members, nil
}

// MemberExists 检查成员是否存在
func (s *memberService) MemberExists(id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}

	member, err := s.memberDao.GetMemberByID(id)
	if err != nil {
		return false, err
	}

	return member != nil, nil
}

// validateMember 验证成员数据
func (s *memberService) validateMember(member *model.Member) error {
	// 验证名称
	if strings.TrimSpace(member.Name) == "" {
		return errors.New("成员名称不能为空")
	}

	if len(member.Name) > 50 {
		return errors.New("成员名称长度不能超过50个字符")
	}

	// 验证邮箱格式（如果提供了邮箱）
	if member.Email != "" {
		if !strings.Contains(member.Email, "@") {
			return errors.New("邮箱格式不正确")
		}

		if len(member.Email) > 100 {
			return errors.New("邮箱长度不能超过100个字符")
		}
	}

	// 验证手机号格式（如果提供了手机号）
	if member.Phone != "" {
		if len(member.Phone) > 20 {
			return errors.New("手机号长度不能超过20个字符")
		}
	}

	// 验证角色
	if !s.isValidRole(member.Role) {
		return errors.New("无效的成员角色")
	}

	return nil
}

// emailExists 检查邮箱是否已被使用
func (s *memberService) emailExists(email string, excludeID uint) (bool, error) {
	members, err := s.memberDao.GetAllMembers()
	if err != nil {
		return false, err
	}

	for _, member := range members {
		if member.Email == email && member.ID != excludeID {
			return true, nil
		}
	}

	return false, nil
}

// phoneExists 检查手机号是否已被使用
func (s *memberService) phoneExists(phone string, excludeID uint) (bool, error) {
	members, err := s.memberDao.GetAllMembers()
	if err != nil {
		return false, err
	}

	for _, member := range members {
		if member.Phone == phone && member.ID != excludeID {
			return true, nil
		}
	}

	return false, nil
}

// isValidRole 验证角色是否有效
func (s *memberService) isValidRole(role model.MemberRole) bool {
	switch role {
	case model.RoleAdmin, model.RoleMember, model.RoleViewer:
		return true
	default:
		return false
	}
}
