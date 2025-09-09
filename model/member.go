package model

import (
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
)

// MemberDao 成员数据访问对象
type MemberDao struct{}

var (
	memberOnce sync.Once
	memberDao  *MemberDao
)

// NewMemberDaoInstance 返回 MemberDao 单例实例
func NewMemberDaoInstance() *MemberDao {
	memberOnce.Do(func() {
		memberDao = &MemberDao{}
	})
	return memberDao
}

// CreateMember 创建成员
func (MemberDao) CreateMember(member *Member) error {
	if err := database.DB.Create(member).Error; err != nil {
		log.Printf("创建成员失败: %v", err)
		return err
	}
	return nil
}

// GetMemberByID 根据ID获取成员
func (MemberDao) GetMemberByID(id uint) (*Member, error) {
	var member Member
	if err := database.DB.Preload("Family").First(&member, id).Error; err != nil {
		log.Printf("获取成员失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &member, nil
}

// GetMembersByFamilyID 根据家庭ID获取成员列表
func (MemberDao) GetMembersByFamilyID(familyID uint) ([]Member, error) {
	var members []Member
	if err := database.DB.Where("family_id = ? AND status = 1", familyID).Find(&members).Error; err != nil {
		log.Printf("获取家庭成员失败 FamilyID=%d: %v", familyID, err)
		return nil, err
	}
	return members, nil
}

// UpdateMember 更新成员信息
func (MemberDao) UpdateMember(member *Member) error {
	if err := database.DB.Model(member).Select("name", "phone", "email").Updates(member).Error; err != nil {
		log.Printf("更新成员失败 ID=%d: %v", member.ID, err)
		return err
	}
	return nil
}

// DeleteMember 软删除成员（设置状态为0）
func (MemberDao) DeleteMember(id uint) error {
	if err := database.DB.Model(&Member{}).Where("id = ?", id).Update("status", 0).Error; err != nil {
		log.Printf("删除成员失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}

func (MemberDao) GetAllMembers() ([]Member, error) {
	var members []Member
	if err := database.DB.Preload("Family").Find(&members).Error; err != nil {
		log.Printf("获取所有成员失败: %v", err)
		return nil, err
	}
	return members, nil
}

func (MemberDao) GetActiveMembersByFamilyID(familyID uint) ([]Member, error) {
	var members []Member
	if err := database.DB.Where("family_id = ? AND status = 1", familyID).Find(&members).Error; err != nil {
		log.Printf("获取家庭活跃成员失败 FamilyID=%d: %v", familyID, err)
		return nil, err
	}
	return members, nil
}

// UpdateMemberRole 更新成员角色
func (MemberDao) UpdateMemberRole(id uint, role MemberRole) error {
	if err := database.DB.Model(&Member{}).Where("id = ?", id).Update("role", role).Error; err != nil {
		log.Printf("更新成员角色失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}
