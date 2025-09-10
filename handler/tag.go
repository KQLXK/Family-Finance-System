// handler/tag_handler.go
package handler

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"github.com/KQLXK/Family-Finance-System/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: service.NewTagService(),
	}
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.tagService.CreateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "标签创建成功",
		"data":    tag,
	})
}

// GetTagByID 根据ID获取标签
func (h *TagHandler) GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	tag, err := h.tagService.GetTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tag,
	})
}

// GetTagsByFamilyID 根据家庭ID获取标签列表
func (h *TagHandler) GetTagsByFamilyID(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	tags, err := h.tagService.GetTagsByFamilyID(uint(familyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tags,
	})
}

// GetTagsByType 根据类型获取标签列表
func (h *TagHandler) GetTagsByType(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	tagType := c.Query("type")
	if tagType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签类型不能为空"})
		return
	}

	tags, err := h.tagService.GetTagsByType(uint(familyID), tagType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tags,
	})
}

// GetAllTags 获取所有标签
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tags,
	})
}

// UpdateTag 更新标签信息
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	tag.ID = uint(id)
	if err := h.tagService.UpdateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签信息更新成功",
		"data":    tag,
	})
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	if err := h.tagService.DeleteTag(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签删除成功",
	})
}
