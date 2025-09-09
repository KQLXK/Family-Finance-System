// handler/member_handler.go
package handler

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"github.com/KQLXK/Family-Finance-System/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MemberHandler 成员处理器
type MemberHandler struct {
	memberService service.MemberService
}

// NewMemberHandler 创建成员处理器
func NewMemberHandler() *MemberHandler {
	return &MemberHandler{
		memberService: service.NewMemberService(),
	}
}

// CreateMember 创建成员
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var member model.Member

	familyIDstr := c.Param("id")
	familyID, err := strconv.ParseInt(familyIDstr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}
	member.FamilyID = uint(familyID)

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.memberService.CreateMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "成员创建成功",
		"data":    member,
	})
}

// GetMemberByID 根据ID获取成员
func (h *MemberHandler) GetMemberByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的成员ID"})
		return
	}

	member, err := h.memberService.GetMemberByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": member,
	})
}

// GetMembersByFamilyID 根据家庭ID获取成员列表
func (h *MemberHandler) GetMembersByFamilyID(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	members, err := h.memberService.GetMembersByFamilyID(uint(familyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}

// GetAllMembers 获取所有成员
func (h *MemberHandler) GetAllMembers(c *gin.Context) {
	members, err := h.memberService.GetAllMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}

// UpdateMember 更新成员信息
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的成员ID"})
		return
	}

	var member model.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	member.ID = uint(id)
	if err := h.memberService.UpdateMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成员信息更新成功",
		"data":    member,
	})
}

// DeleteMember 删除成员
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的成员ID"})
		return
	}

	if err := h.memberService.DeleteMember(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成员删除成功",
	})
}

// ChangeMemberRole 更改成员角色
func (h *MemberHandler) ChangeMemberRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的成员ID"})
		return
	}

	var request struct {
		Role model.MemberRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.memberService.ChangeMemberRole(uint(id), request.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成员角色更新成功",
	})
}

// GetActiveMembersByFamilyID 根据家庭ID获取活跃成员列表
func (h *MemberHandler) GetActiveMembersByFamilyID(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	members, err := h.memberService.GetActiveMembersByFamilyID(uint(familyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}
