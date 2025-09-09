// handler/family_handler.go
package handler

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"github.com/KQLXK/Family-Finance-System/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FamilyHandler 家庭处理器
type FamilyHandler struct {
	familyService service.FamilyService
}

// NewFamilyHandler 创建家庭处理器
func NewFamilyHandler() *FamilyHandler {
	return &FamilyHandler{
		familyService: service.NewFamilyService(),
	}
}

// CreateFamily 创建家庭
func (h *FamilyHandler) CreateFamily(c *gin.Context) {
	var family model.Family
	if err := c.ShouldBindJSON(&family); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if family.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "家庭名称不能为空"})
		return
	}

	if err := h.familyService.CreateFamily(&family); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建家庭失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "家庭创建成功",
		"data":    family,
	})
}

// GetFamilyByID 根据ID获取家庭
func (h *FamilyHandler) GetFamilyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	family, err := h.familyService.GetFamilyByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "家庭不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": family,
	})
}

// GetAllFamilies 获取所有家庭
func (h *FamilyHandler) GetAllFamilies(c *gin.Context) {
	families, err := h.familyService.GetAllFamilies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取家庭列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": families,
	})
}

// UpdateFamily 更新家庭信息
func (h *FamilyHandler) UpdateFamily(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	var family model.Family
	if err := c.ShouldBindJSON(&family); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	family.ID = uint(id)
	if err := h.familyService.UpdateFamily(&family); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新家庭信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "家庭信息更新成功",
		"data":    family,
	})
}

// DeleteFamily 删除家庭
func (h *FamilyHandler) DeleteFamily(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	if err := h.familyService.DeleteFamily(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除家庭失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "家庭删除成功",
	})
}
